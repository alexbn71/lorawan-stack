with import <nixpkgs> {};

let
  userName = "admin";
  appName = "test-app";
  devName = "test-dev";

  frequencyPlan = "EU_863_870";
  #frequencyPlan = "US_902_928_FSB_1";

  stackPort = "1885";
  udpPort = "1700";

  cliCmd = "${go}/bin/go run ./cmd/ttn-lw-cli";
  stackCmd = "${go}/bin/go run ./cmd/ttn-lw-stack";

  createAPIKeyCmd = "${cliCmd} application api-keys create --application-id ${appName} --right-application-all | ${jq}/bin/jq -c -r '.key'";

  makeWaitCmd = event: ''
    read -n 1 -p "Ensure that ${event} and press any key to continue"
  '';

  joinWaitCmd = makeWaitCmd "the device has joined";
  bootWaitCmd = makeWaitCmd "the device has booted";
  uplinkWaitCmd = makeWaitCmd "the device has sent an uplink";

  linkApp = writeShellScriptBin "link-app" ''
    set -xe
    apiKey="$(${createAPIKeyCmd})"
    ${cliCmd} application link set ${appName} --api-key ''${apiKey} ''${@}
    printf "export TTN_API_KEY=\"''${apiKey}\"" > .envrc.api_key
    ${direnv}/bin/direnv reload
  '';

  bootstrap = writeShellScriptBin "bootstrap" ''
    set -xe
    ${cliCmd} login
    ${cliCmd} application create ${appName} --user-id ${userName}
    printf "export BOOTSTRAP_DONE=1" > .envrc.bootstrap
    ${linkApp}/bin/link-app
  '';

  scheduleClockExample = writeShellScriptBin "schedule-clock-example" ''
    set -xe
    cd ../../github.com/TheThingsNetwork/lorawan-stack-example-clock/controller
    export TTN_APPLICATION_ID=''${TTN_APPLICATION_ID:-${appName}}
    export TTN_DEVICE_ID=''${TTN_DEVICE_ID:-${devName}}
    export TTN_API_KEY=''${TTN_API_KEY:-$(cat .envrc.api_key)}
    export TTN_SERVER=''${TTN_SERVER:-http://localhost:${stackPort}}
    ${nodejs}/bin/npm run start
  '';

  # TODO: Support class B/C args
  makePushDownlinksCmd = lib.concatMapStringsSep "\n" (down: ''
     ${cliCmd} device downlink push ${appName} ${devName} \
                                    --f-port=${toString down.port} \
                                    --frm-payload=${down.payload}
      sleep 1
  '');

  classADownlinks = [
    {
      port = 1;
      payload = "1111";
    }
    {
      port = 42;
      payload = "22";
    }
    {
      port = 120;
      payload = "333333";
    }
  ];

  pushClassADownlinksCmd = makePushDownlinksCmd classADownlinks;
  pushClassCDownlinksCmd = makePushDownlinksCmd (classADownlinks ++ [
    # TODO: Add class C downlinks
  ]);
in
  stdenv.mkDerivation {
    CGO_ENABLED = "0";
    GO111MODULE = "on";
    GOROOT = "${go}/share/go";
    TTN_LW_GS_UDP_LISTENERS=":${udpPort}=${frequencyPlan}";

    name = "ttn-env";
    buildInputs = [
      coreutils
      gnumake
      gnused
      go
      go-tools
      gotools
      jq
      minicom
      mosquitto
      neovim
      nodejs
      perl
      protobuf
      redis
      richgo
      rpm
      travis

      (writeShellScriptBin "yarn"
      ''
        ${toString ./node_modules/.bin/yarn} ''${@}
      '')
      (writeShellScriptBin "mage"
      ''
        ${toString ./mage} ''${@}
      '')

      (writeShellScriptBin "cli"
      ''
        ${cliCmd} ''${@}
      '')
      (writeShellScriptBin "stack"
      ''
        ${stackCmd} ''${@}
      '')

      (writeShellScriptBin "snapcraft"
      ''
        ${docker}/bin/docker run --rm -i \
          --user $(id -u):$(id -g) \
          --mount type=bind,src=$XDG_CONFIG_HOME/snapcraft,dst=$XDG_CONFIG_HOME/snapcraft \
          --mount type=bind,src=$(pwd),dst=$(pwd) \
          --entrypoint "/snap/bin/snapcraft" \
          -e XDG_CONFIG_HOME \
          -w $(pwd) \
          snapcore/snapcraft:candidate ''${@}
      '')

      (writeShellScriptBin "start-clean-stack"
      ''
        set -xe
        ${gnumake}/bin/make dev.databases.stop
        sudo ${coreutils}/bin/rm -rf .dev/data
        ${gnumake}/bin/make dev.stack.init
        printf "export BOOTSTRAP_DONE=0" > .envrc.bootstrap
        ${direnv}/bin/direnv reload
        ${stackCmd} start ''${@}
      '')

      (writeShellScriptBin "create-app" ''
        set -xe
        ${cliCmd} application create ${appName} --user-id ${userName} ''${@}
      '')
      (writeShellScriptBin "delete-app" ''
        set -xe
        ${cliCmd} application delete ${appName} ''${@}
      '')

      (writeShellScriptBin "create-dev" ''
        set -xe
        ${cliCmd} device create ${appName} ${devName} ''${@}
      '')
      (writeShellScriptBin "delete-dev" ''
        set -xe
        ${cliCmd} device delete ${appName} ${devName} ''${@}
      '')

      (writeShellScriptBin "push-down" ''
        set -xe
        ${cliCmd} device downlink push ${appName} ${devName} ''${@}
      '')


      bootstrap
      linkApp
      scheduleClockExample

      # TODO: Merge into test-device
      (writeShellScriptBin "test-class-a-device-otaa" ''
        case "''${1}" in
        "ff1705")
          devEUI="0102030405060708"
          joinEUI="0102030405060708"
          appKey="01020304050607080102030405060708"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        "ff1705-semtech")
          devEUI="ffffffaa11111110"
          joinEUI="0102030405060708"
          appKey="''${SEMTECH_APPKEY}"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        "simulator")
          devEUI="4200000000000000"
          joinEUI="4200000000000000"
          appKey="42000000000000000000000000000000"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        "node")
          devEUI="0004A30B001BC7AF"
          joinEUI="4200000000000000"
          appKey="01020304050607080102030405060708"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        "uno")
          devEUI="0102030405060708"
          joinEUI="0102030405060708"
          appKey="01020304050607080102030405060708"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        *)
          echo "Unknown device ''${1} - specify either 'uno' or 'simulator'"
          exit 1
        esac
        shift

        [ "''${BOOTSTRAP_DONE}" == "1" ] || ${bootstrap}/bin/bootstrap

        set -xe
        ${cliCmd} device delete ${appName} ${devName} || true
        ${cliCmd} device create ${appName} ${devName} \
                        --join-eui ''${joinEUI} \
                        --dev-eui ''${devEUI} \
                        --root-keys.app-key.key ''${appKey} \
                        --lorawan_phy_version ''${phyVersion} \
                        --lorawan-version ''${macVersion} \
                        --frequency_plan_id ${frequencyPlan} \
                        ''${@}

        ${joinWaitCmd}

        ${cliCmd} device get ${appName} ${devName}

        ${makePushDownlinksCmd [{
          port = 2;
          payload = "abcd";
        }]}

        ${uplinkWaitCmd}

        ${pushClassADownlinksCmd}
      '')
      (writeShellScriptBin "test-class-a-device-abp" ''
        extra=""
        case "''${1}" in
        "node")
          devEUI="0004A30B001BC7AF"
          joinEUI="4200000000000000"
          devAddr="42FFFFFF"
          appSKey="42000000000000000000000000000000"
          nwkSKey="42000000000000000000000000000042"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          extra="--mac-settings.factory-preset-frequencies 868100000,868300000,868500000,867100000,867300000,867500000,867700000,867900000 \
                 --mac-settings.resets-f-cnt=true \
                 --mac-settings.rx2-data-rate-index DATA_RATE_3"
          ;;
        *)
          echo "Unknown device ''${1} - specify either 'node' or ?"
          exit 1
        esac
        shift

        [ "''${BOOTSTRAP_DONE}" == "1" ] || ${bootstrap}/bin/bootstrap

        set -xe
        ${cliCmd} device delete ${appName} ${devName} || true
        ${cliCmd} device create ${appName} ${devName} \
                        --join-eui ''${joinEUI} \
                        --dev-eui ''${devEUI} \
                        --supports-join false \
                        --session.dev_addr ''${devAddr} \
                        --session.keys.app_s_key.key ''${appSKey} \
                        --session.keys.f_nwk_s_int_key.key ''${nwkSKey} \
                        --lorawan_phy_version ''${phyVersion} \
                        --lorawan-version ''${macVersion} \
                        --frequency_plan_id ${frequencyPlan} \
                        ''${extra} \
                        ''${@}

        ${bootWaitCmd}

        ${cliCmd} device get ${appName} ${devName}

        ${makePushDownlinksCmd [{
          port = 2;
          payload = "abcd";
        }]}

        ${uplinkWaitCmd}

        ${pushClassADownlinksCmd}
      '')
      (writeShellScriptBin "test-class-c-device"
      ''
        case "''${1}" in
        "ff1705")
          devEUI="0102030405060708"
          joinEUI="0102030405060708"
          appKey="01020304050607080102030405060708"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        "ff1705-semtech")
          devEUI="ffffffaa11111110"
          joinEUI="0102030405060708"
          appKey="''${SEMTECH_APPKEY}"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        "simulator")
          devEUI="4200000000000000"
          joinEUI="4200000000000000"
          appKey="42000000000000000000000000000000"
          macVersion="1.0.2"
          phyVersion="1.0.2-b"
          ;;
        *)
          echo "Unknown device ''${1} - specify either 'ff1705' or 'simulator'"
          exit 1
        esac
        shift

        [ "''${BOOTSTRAP_DONE}" == "1" ] || ${bootstrap}/bin/bootstrap

        set -xe
        ${cliCmd} device delete ${appName} ${devName} || true
        ${cliCmd} device create ${appName} ${devName} \
                        --join-eui ''${joinEUI} \
                        --dev-eui ''${devEUI} \
                        --root-keys.app-key.key ''${appKey} \
                        --lorawan_phy_version ''${phyVersion} \
                        --lorawan-version ''${macVersion} \
                        --frequency_plan_id ${frequencyPlan} \
                        --supports-class-c \
                        --resets-join-nonces \
                        ''${@}

        ${joinWaitCmd}

        ${cliCmd} device get ${appName} ${devName}

        source .envrc.api_key
        ${scheduleClockExample}/bin/schedule-clock-example

        ${pushClassCDownlinksCmd}
      '')
    ];
  }
