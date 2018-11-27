// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package identityserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/pkg/auth"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/pkg/rpcmetadata"
	"go.thethings.network/lorawan-stack/pkg/rpcmiddleware/warning"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
)

var (
	errUnauthenticated          = errors.DefineUnauthenticated("unauthenticated", "unauthenticated")
	errUnsupportedAuthorization = errors.DefineUnauthenticated("unsupported_authorization", "Unsupported authorization method")
	errInvalidAuthorization     = errors.DefinePermissionDenied("invalid_authorization", "invalid authorization")
	errTokenExpired             = errors.DefineUnauthenticated("token_expired", "access token expired")
	errOAuthClientRejected      = errors.DefinePermissionDenied("oauth_client_rejected", "OAuth client was rejected")
	errOAuthClientSuspended     = errors.DefinePermissionDenied("oauth_client_suspended", "OAuth client was suspended")
)

type requestAccessKeyType struct{}

var requestAccessKey requestAccessKeyType

type requestAccess struct {
	authInfo     *ttnpb.AuthInfoResponse
	entityRights map[*ttnpb.EntityIdentifiers]*ttnpb.Rights
}

func (is *IdentityServer) withRequestAccessCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestAccessKey, new(requestAccess))
}

func (is *IdentityServer) authInfo(ctx context.Context) (info *ttnpb.AuthInfoResponse, err error) {
	if access, ok := ctx.Value(requestAccessKey).(*requestAccess); ok {
		if access.authInfo != nil {
			return access.authInfo, nil
		}
		defer func() {
			if err == nil {
				access.authInfo = info
			}
		}()
	}

	md := rpcmetadata.FromIncomingContext(ctx)
	if md.AuthType == "" {
		return &ttnpb.AuthInfoResponse{}, nil
	}
	if strings.ToLower(md.AuthType) != "bearer" {
		return nil, errUnsupportedAuthorization
	}
	token := md.AuthValue

	if info = is.cachedAuthInfoForToken(ctx, token); info != nil {
		return info, nil
	}
	defer func() {
		if err == nil {
			is.cacheAuthInfoForToken(ctx, token, info)
		}
	}()

	tokenType, tokenID, tokenKey, err := auth.SplitToken(token)
	if err != nil {
		return nil, err
	}

	var fetch func(db *gorm.DB) error
	res := new(ttnpb.AuthInfoResponse)
	userFieldMask := &types.FieldMask{Paths: []string{"admin", "state"}}
	clientFieldMask := &types.FieldMask{Paths: []string{"state"}}
	var user *ttnpb.User
	var userRights *ttnpb.Rights

	switch tokenType {
	case auth.APIKey:
		fetch = func(db *gorm.DB) error {
			ids, apiKey, err := store.GetAPIKeyStore(db).GetAPIKey(ctx, tokenID)
			if err != nil {
				return err
			}
			valid, err := auth.Password(apiKey.GetKey()).Validate(tokenKey)
			if err != nil {
				return err
			}
			if !valid {
				return errInvalidAuthorization
			}
			apiKey.Key = ""
			res.AccessMethod = &ttnpb.AuthInfoResponse_APIKey{
				APIKey: &ttnpb.AuthInfoResponse_APIKeyAccess{
					APIKey:    *apiKey,
					EntityIDs: *ids,
				},
			}
			if userIDs := ids.GetUserIDs(); userIDs != nil {
				user, err = store.GetUserStore(db).GetUser(ctx, userIDs, userFieldMask)
				if err != nil {
					return err
				}
				userRights = ttnpb.RightsFrom(apiKey.Rights...)
			}
			return nil
		}
	case auth.AccessToken:
		fetch = func(db *gorm.DB) error {
			accessToken, err := store.GetOAuthStore(db).GetAccessToken(ctx, tokenID)
			if err != nil {
				return err
			}
			valid, err := auth.Password(accessToken.GetAccessToken()).Validate(tokenKey)
			if err != nil {
				return err
			}
			if !valid {
				return errInvalidAuthorization
			}
			if accessToken.ExpiresAt.Before(time.Now()) {
				return errTokenExpired
			}
			accessToken.AccessToken, accessToken.RefreshToken = "", ""
			res.AccessMethod = &ttnpb.AuthInfoResponse_OAuthAccessToken{
				OAuthAccessToken: accessToken,
			}
			user, err = store.GetUserStore(db).GetUser(ctx, &accessToken.UserIDs, userFieldMask)
			if err != nil {
				return err
			}
			client, err := store.GetClientStore(db).GetClient(ctx, &accessToken.ClientIDs, clientFieldMask)
			if err != nil {
				return err
			}
			switch client.State {
			case ttnpb.STATE_REQUESTED:
				// OAuth authorization only passes for collaborators, so this is ok.
			case ttnpb.STATE_APPROVED:
				// Normal OAuth client.
			case ttnpb.STATE_REJECTED:
				return errOAuthClientRejected
			case ttnpb.STATE_FLAGGED:
				// Innocent until proven guilty.
			case ttnpb.STATE_SUSPENDED:
				return errOAuthClientSuspended
			default:
				panic(fmt.Sprintf("Unhandled client state: %s", client.State.String()))
			}
			userRights = ttnpb.RightsFrom(accessToken.Rights...)
			return nil
		}
	default:
		return nil, errUnsupportedAuthorization
	}

	err = is.withDatabase(ctx, fetch)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if user.Admin {
			res.UniversalRights = ttnpb.AllRights.Implied().Intersect(userRights.Implied()) // TODO: Use restricted Admin rights.
		}

		switch user.State {
		case ttnpb.STATE_REQUESTED:
			// Go to profile page, edit basic settings (such as email), delete account.
			restrictRights(res, ttnpb.RightsFrom(ttnpb.RIGHT_USER_INFO, ttnpb.RIGHT_USER_SETTINGS_BASIC, ttnpb.RIGHT_USER_DELETE))
			warning.Add(ctx, "Restricted rights while account pending")
		case ttnpb.STATE_APPROVED:
			// Normal user.
		case ttnpb.STATE_REJECTED:
			// Go to profile page, delete account.
			restrictRights(res, ttnpb.RightsFrom(ttnpb.RIGHT_USER_INFO, ttnpb.RIGHT_USER_DELETE))
			warning.Add(ctx, "Restricted rights after account rejection")
		case ttnpb.STATE_FLAGGED:
			// Innocent until proven guilty.
		case ttnpb.STATE_SUSPENDED:
			// Go to profile page.
			restrictRights(res, ttnpb.RightsFrom(ttnpb.RIGHT_USER_INFO))
			warning.Add(ctx, "Restricted rights after account suspension")
		default:
			panic(fmt.Sprintf("Unhandled user state: %s", user.State.String()))
		}
	}

	return res, nil
}

func (is *IdentityServer) RequireAuthenticated(ctx context.Context) error {
	authInfo, err := is.authInfo(ctx)
	if err != nil {
		return err
	}
	if apiKey := authInfo.GetAPIKey(); apiKey != nil {
		return nil
	} else if accessToken := authInfo.GetOAuthAccessToken(); accessToken != nil {
		return nil
	}
	return errUnauthenticated
}

func (is *IdentityServer) UniversalRights(ctx context.Context) *ttnpb.Rights {
	info, err := is.authInfo(ctx)
	if err == nil {
		return info.GetUniversalRights()
	}
	return nil
}

func restrictRights(info *ttnpb.AuthInfoResponse, rights *ttnpb.Rights) {
	if apiKey := info.GetAPIKey(); apiKey != nil {
		apiKey.Rights = ttnpb.RightsFrom(apiKey.Rights...).Intersect(rights).GetRights()
	} else if token := info.GetOAuthAccessToken(); token != nil {
		token.Rights = ttnpb.RightsFrom(token.Rights...).Intersect(rights).GetRights()
	}
	info.UniversalRights = info.UniversalRights.Intersect(rights)
}

func entityRights(authInfo *ttnpb.AuthInfoResponse) (*ttnpb.EntityIdentifiers, *ttnpb.Rights) {
	if apiKey := authInfo.GetAPIKey(); apiKey != nil {
		return &apiKey.EntityIDs, ttnpb.RightsFrom(apiKey.Rights...)
	} else if accessToken := authInfo.GetOAuthAccessToken(); accessToken != nil {
		return accessToken.UserIDs.EntityIdentifiers(), ttnpb.RightsFrom(accessToken.Rights...)
	}
	return nil, nil
}

func (is *IdentityServer) entityRights(ctx context.Context, authInfo *ttnpb.AuthInfoResponse) (res map[*ttnpb.EntityIdentifiers]*ttnpb.Rights, err error) {
	if access, ok := ctx.Value(requestAccessKey).(*requestAccess); ok {
		if access.entityRights != nil {
			return access.entityRights, nil
		}
		defer func() {
			if err == nil {
				access.entityRights = res
			}
		}()
	}

	ids, rights := entityRights(authInfo)
	if ids == nil {
		return nil, nil
	}
	entityRights := make(map[*ttnpb.EntityIdentifiers]*ttnpb.Rights)
	entityRights[ids] = rights.Implied()
	memberRights, err := is.memberRights(ctx, ids)
	if err != nil {
		return nil, err
	}
	for ids, memberRights := range memberRights {
		entityRights[ids] = memberRights.Implied().Intersect(rights.Implied())
	}
	return entityRights, nil
}

func (is *IdentityServer) memberRights(ctx context.Context, ids *ttnpb.EntityIdentifiers) (entityRights map[*ttnpb.EntityIdentifiers]*ttnpb.Rights, err error) {
	var ouIDs *ttnpb.OrganizationOrUserIdentifiers
	switch ids := ids.Identifiers().(type) {
	case *ttnpb.OrganizationIdentifiers:
		ouIDs = ids.OrganizationOrUserIdentifiers()
	case *ttnpb.UserIdentifiers:
		ouIDs = ids.OrganizationOrUserIdentifiers()
	}
	if ouIDs == nil {
		return nil, nil
	}

	memberships := is.cachedMembershipsForAccount(ctx, ouIDs)
	if memberships == nil {
		err = is.withDatabase(ctx, func(db *gorm.DB) error {
			memberships, err = store.GetMembershipStore(db).FindMemberRights(ctx, ouIDs, "")
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		defer func() {
			is.cacheMembershipsForAccount(ctx, ouIDs, memberships)
		}()
	}

	entityRights = make(map[*ttnpb.EntityIdentifiers]*ttnpb.Rights)
	for ids, rights := range memberships {
		entityRights[ids] = rights
		subMemberRights, err := is.memberRights(ctx, ids)
		if err != nil {
			return nil, err
		}
		for ids, memberRights := range subMemberRights {
			entityRights[ids] = memberRights.Implied().Intersect(rights.Implied())
		}
	}
	return entityRights, nil
}

type entityAccess struct {
	*IdentityServer
}

func (ea *entityAccess) AuthInfo(ctx context.Context, _ *types.Empty) (*ttnpb.AuthInfoResponse, error) {
	return ea.authInfo(ctx)
}