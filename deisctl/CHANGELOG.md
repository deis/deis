### initial commit -> v0.12.0

#### Features

 - [`938f4eb`](https://github.com/deis/deis/deisctl/commit/938f4eb3c2c490d0f92dce25d1fd4f684012dfee) cmd: add informative messages to install
 - [`c45fade`](https://github.com/deis/deis/deisctl/commit/c45fadeebe09026599d90193ffd594016cd7e7e7) deisctl: colorize deisctl
 - [`615a2cb`](https://github.com/deis/deis/deisctl/commit/615a2cb40aa7585428f596114b4270027e1e8ab5) start/stop: allow starting or stopping > 1 unit at a time
 - [`077429c`](https://github.com/deis/deis/deisctl/commit/077429cb472553ef2818c9fc27653d574203f3e7) start: wait on containers to start
 - [`74e9337`](https://github.com/deis/deis/deisctl/commit/74e9337f174a928d68d33f657afa4655bce2582b) installer: use /usr/local/bin and one-liner install scripts
 - [`2ae953e`](https://github.com/deis/deis/deisctl/commit/2ae953ece143a84fe891f10d039902acef28945b) Makefile: create shell script installer
 - [`6d541aa`](https://github.com/deis/deis/deisctl/commit/6d541aae5dd141ea8ee3bd1434c6d233770fb37d) config: first pass at config subcommand
 - [`c182966`](https://github.com/deis/deis/deisctl/commit/c1829662a1f02eadf7902963f62207f22ffaf5f4) restart: add restart command for convenience
 - [`bf68a8e`](https://github.com/deis/deis/deisctl/commit/bf68a8e64d9095ddbc151ca61fae5a868646aff7) journal: add journal support
 - [`65ee91e`](https://github.com/deis/deis/deisctl/commit/65ee91eaafdebf672c78d3205e7354806ecb8922) deisctl: move server env variable to etcd
 - [`975ac6b`](https://github.com/deis/deis/deisctl/commit/975ac6bb840bd492d1863564691598ed91233868) deisctl:  add new feautes and update core-os updatectl to updateservicectl
 - [`b366679`](https://github.com/deis/deis/deisctl/commit/b366679d4111b6b098f791ae9f559ed999c819e3) deisctl: get groupid and app id from etcd defaults to env variables
 - [`8509f01`](https://github.com/deis/deis/deisctl/commit/8509f0147dcc8921b781ec438876eadbfa1b8783) deisctl: removed unneccessary code
 - [`558bfb0`](https://github.com/deis/deis/deisctl/commit/558bfb0660ce1eea91625dfa8ae51f05bd928afe) deisctl:  add hooks and units dirs to constants
 - [`cdd0150`](https://github.com/deis/deis/deisctl/commit/cdd015093622f30e50f34aa3837892c0ec3433b1) hook: add pre/post update hooks
 - [`74f5229`](https://github.com/deis/deis/deisctl/commit/74f5229658ad61a7583e166d04cca63bccf75078) deisctl:  working version of updater latest
 - [`d96d9bb`](https://github.com/deis/deis/deisctl/commit/d96d9bbc1b97bd1c58c1bd0d9b9855ccf9453ee7) deisctl:  basic working version of updater
 - [`4c509b1`](https://github.com/deis/deis/deisctl/commit/4c509b1055e629aab3eb743688b8da90d5a32fb0) deisctl:  change flag to string
 - [`5a5f859`](https://github.com/deis/deis/deisctl/commit/5a5f85988b9b28a9b123b6bee611d36cf6db66b2) deisctl: fix package names
 - [`d07a923`](https://github.com/deis/deis/deisctl/commit/d07a923d167535380441f3c1b7e86d2502cdfcf1) deisctl: updated command instance
 - [`cb18218`](https://github.com/deis/deis/deisctl/commit/cb18218d3661ae94854166df759e59853273cd1b) deisctl: removed systemd and distibuted lock
 - [`82ccd66`](https://github.com/deis/deis/deisctl/commit/82ccd6689aa9180e5d075840d64ba2068efb4803) deisctl: add utils for client instance functions
 - [`dcd22f1`](https://github.com/deis/deis/deisctl/commit/dcd22f1b14571b9b65bfb4a7a9e11ec0b31bd19f) deisctl:  add update command and updatectl package

#### Fixes

 - [`be4bef3`](https://github.com/deis/deis/deisctl/commit/be4bef39543eb5bc391899e18ac95b65627a46f7) state: ignore intermittent timeouts when polling for UnitState
 - [`e316572`](https://github.com/deis/deis/deisctl/commit/e31657238b1e5a5f5b6d82e821f1300936370aa0) Makefile: ensure the installer makes /var/lib/deis/units readable
 - [`d831066`](https://github.com/deis/deis/deisctl/commit/d8310669d884a7a1b7cb152415c74001229f593b) registry: fix duplicate tag in registry-data unit file
 - [`561fb31`](https://github.com/deis/deis/deisctl/commit/561fb31baa434881886cb810074ddecf43c62153) units: show start-pre status when downloading data container base
 - [`81371d5`](https://github.com/deis/deis/deisctl/commit/81371d59c20325051e7287326a6b9648e53d2a61) cmd: only allow the router to scale past 1
 - [`05f027c`](https://github.com/deis/deis/deisctl/commit/05f027c6fc89ebc731da91ead64628c0e7e093da) cmd: create unit files as readable to all users
 - [`043739e`](https://github.com/deis/deis/deisctl/commit/043739ea7f87a70cf087cf17d98ecb9c02367e09) client: destroy all units if none specified
 - [`40b084f`](https://github.com/deis/deis/deisctl/commit/40b084f4384b83a79061dfc2a920649b2a0581fa) client: dramatically simplify scaling logic
 - [`37c348e`](https://github.com/deis/deis/deisctl/commit/37c348ef2f94ef62b51b09381c41f9abccea85e6) cmd: append "@1" if none supplied to install
 - [`8f14e8a`](https://github.com/deis/deis/deisctl/commit/8f14e8aea22f415b5e14624a4371b15b5baed327) client: destroy all units if none specified
 - [`725ae0d`](https://github.com/deis/deis/deisctl/commit/725ae0d84e51e466f31cd864521a856773015011) client: check if unit exists
 - [`7bd1709`](https://github.com/deis/deis/deisctl/commit/7bd17097de3e1a6d0f98d86daa95d863323d69ef) client: return error if unit list is empty
 - [`1f79c8b`](https://github.com/deis/deis/deisctl/commit/1f79c8bbc6f09861b2af05a20c31745430d0e92b) deisctl: use docopt's native version parser
 - [`fd18f5b`](https://github.com/deis/deis/deisctl/commit/fd18f5b77491a5950ddc56966f33a2f8eb71c308) makefile: expand paths for golint
 - [`cc9e0a4`](https://github.com/deis/deis/deisctl/commit/cc9e0a49785edf0b11be97d3b75df44464415570) units: add fleet.sock bind-mount for controller
 - [`4497af2`](https://github.com/deis/deis/deisctl/commit/4497af28bb2de71205623a25462fd374b2eb6180) cmd: allow `-p` to specify where to save local unit files
 - [`318514e`](https://github.com/deis/deis/deisctl/commit/318514e3bbf8f11577e70b65a52a180dbcab7691) client: look for unit files in ~/.deisctl before /var/lib/deis/units
 - [`c325ca2`](https://github.com/deis/deis/deisctl/commit/c325ca28c085acc5ffbec8432cec31d8e820da8e) debug: remove other vestige of unused --debug flag
 - [`72cb9d8`](https://github.com/deis/deis/deisctl/commit/72cb9d89fddded0057fc6de416170236eef45163) README: fix installer link to use http, not https
 - [`54dc7df`](https://github.com/deis/deis/deisctl/commit/54dc7dfc9aaea32a1f2d16897d05c90f2e56171c) README: update installer link
 - [`ba85174`](https://github.com/deis/deis/deisctl/commit/ba851748bde851f61462a0010d7c5f0d7055e155) installer: use deisctl-hack fork of makeself
 - [`dce122d`](https://github.com/deis/deis/deisctl/commit/dce122d2ec7a770b6f3dcb168ee591a83f6e6bc6) debug: remove unused --debug option
 - [`9459a83`](https://github.com/deis/deis/deisctl/commit/9459a835089a6fa7c3fd37d5e3d8b4da22214033) version: add special handling for --version
 - [`3aaf764`](https://github.com/deis/deis/deisctl/commit/3aaf764cceadb29d88499d716f669e38ef03359c) cmd: add explicit platform target
 - [`4b9f157`](https://github.com/deis/deis/deisctl/commit/4b9f157dbd63b0b3589b7f4bd3d6174b851aab8a) tests: explicitly set tunnel to null
 - [`dce34ce`](https://github.com/deis/deis/deisctl/commit/dce34ce3edc89f6eb12a334a2963625c46e5ab71) destroy: fix shadowing bug in destroy
 - [`3227656`](https://github.com/deis/deis/deisctl/commit/3227656ad05a171b30cb63f1c8b9cdd62f619744) units: controller waits for logger container in ExecStartPre
 - [`471e4e8`](https://github.com/deis/deis/deisctl/commit/471e4e87c899288ea130962c97be68103db18012) units: use @ in wildcard for router conflict
 - [`c2b75ee`](https://github.com/deis/deis/deisctl/commit/c2b75ee318ed1f6adaeadd941908d64b7d89f881) unit: match @ units properly
 - [`f6e0b86`](https://github.com/deis/deis/deisctl/commit/f6e0b86d21a077c8376c5eb17eca4e3a055c004d) destroy: wait for job state inactive on destroy
 - [`e6d260a`](https://github.com/deis/deis/deisctl/commit/e6d260ad150d185c53928b7d20486fb49f4db2f0) ssh: switch to default known hosts
 - [`cba85ee`](https://github.com/deis/deis/deisctl/commit/cba85eeff5797fb9957a22695b5e37c2792801f3) units: use default GOPATH for unit lookup, if available
 - [`a6dc5a9`](https://github.com/deis/deis/deisctl/commit/a6dc5a9818a43ce3895c3d4d00c1f0484494cc12) update: fix imports
 - [`700ad7b`](https://github.com/deis/deis/deisctl/commit/700ad7b1185c39369a4fa0972fe07553df89b23b) state: print inactive states without substates
 - [`265cc3d`](https://github.com/deis/deis/deisctl/commit/265cc3d8c765b42c7eb23dda0e4dc4291ee1df02) units: switch to systemd template units
 - [`985f003`](https://github.com/deis/deis/deisctl/commit/985f0039729bfdf8843344d4040b14668dd9b882) deisctl: fix utils error
 - [`7358f4d`](https://github.com/deis/deis/deisctl/commit/7358f4d7ea7dc28d35427fd7b14769745d9314cf) update: extract update to root
 - [`1d629e5`](https://github.com/deis/deis/deisctl/commit/1d629e566858456dfe6a1b34606ed38ad3358f12) update: add update service as systemd unit
 - [`d6ccce2`](https://github.com/deis/deis/deisctl/commit/d6ccce20c906d1ef2ec07a4c8b8a1149910544bd) updatectl: fix data container matching, fallback to envvar for version
 - [`88745a8`](https://github.com/deis/deis/deisctl/commit/88745a825d66d11caef6d588ae657268bb6f21f1) update: do not pull images on update
 - [`329c372`](https://github.com/deis/deis/deisctl/commit/329c372cd153eb481da25de20803262e56fda32b) constant: add new constant package
 - [`7c2c5db`](https://github.com/deis/deis/deisctl/commit/7c2c5db2c7d42c015e39d1d77ad2b372274a2518) (all): rename constant folder, go fmt
 - [`63e7db9`](https://github.com/deis/deis/deisctl/commit/63e7db9f804491aebeb321b5c751634c27f099e4) units: cleanup post-start output for builder/registry
 - [`910ccef`](https://github.com/deis/deis/deisctl/commit/910ccef6c735a02ff0e773f488969bf21c26aa7d) install: install registry after cache
 - [`7687c1d`](https://github.com/deis/deis/deisctl/commit/7687c1d7e26fe3a2a52041d53eeba7110d1da2e1) units: switch to new fleet X-ConditionMachineID
 - [`4be9b61`](https://github.com/deis/deis/deisctl/commit/4be9b6192b573f55d1fe26515a72911efaac1716) packaging: add version to package tarball
 - [`557cefd`](https://github.com/deis/deis/deisctl/commit/557cefde26489774e19d483b8d76554f611cb5cf) packaging: update Dockerfile and paths
 - [`f7363fa`](https://github.com/deis/deis/deisctl/commit/f7363fae83f1ee9c46c94198cc73c7d1939666bc) upstream: rebase against fleet upstream changes
 - [`5602a17`](https://github.com/deis/deis/deisctl/commit/5602a17f2573c0d4ae3b4dce78bf9dbf99185006) main: fix package path

#### Documentation

 - [`1c30515`](https://github.com/deis/deis/deisctl/commit/1c30515f5aff816167159131ab5fc0cec0c31557) README: update dev documentation
 - [`97b553a`](https://github.com/deis/deis/deisctl/commit/97b553af90a24a3c4494ac3d16a9118fafa9d966) README: link to latest installers on S3, omit "how to build"
 - [`486b422`](https://github.com/deis/deis/deisctl/commit/486b4225db58f1bfd44437e094179c1ec5c7280e) readme: update install instructions
 - [`e65ffcb`](https://github.com/deis/deis/deisctl/commit/e65ffcbf1968ce1acf7105c938c86d2a40971d0d) readme: minor language updates
 - [`b08541d`](https://github.com/deis/deis/deisctl/commit/b08541dd6da1d32c57b153e217f41df8c919ac2e) readme: first pass at readme

#### Maintenance

 - [`1a38aff`](https://github.com/deis/deis/deisctl/commit/1a38aff5d22a7bee965029e54fbdeb4b3a0b8fbb) units: remove deprecated X-Condition from fleet units
 - [`006556e`](https://github.com/deis/deis/deisctl/commit/006556e9ef8c437d0a1204b00c6e4902dd574e92) README: update current version to 0.12.0-dev
 - [`539ed23`](https://github.com/deis/deis/deisctl/commit/539ed23e7d972b6f8bda81b92f618ed09f8c08a3) deictl: bump version in sync with Deis
 - [`23301e8`](https://github.com/deis/deis/deisctl/commit/23301e84b52fb0354cbcfcad451f96627a8c4d53) godeps: bump fleet, updateservicectl, docker
 - [`3d0bf7f`](https://github.com/deis/deis/deisctl/commit/3d0bf7f17d916cf24f15cde34da4707a519862a8) flags: switch to DEISCTL_TUNNEL
 - [`b610417`](https://github.com/deis/deis/deisctl/commit/b6104173b1c6fadcf311650335b556670fcdbcc4) version: bump to 0.11.0
