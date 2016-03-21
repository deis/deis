### v1.12.3 -> v1.13.0

#### Features

 - [`3181b2d`](https://github.com/deis/deis/commit/3181b2d4c70c8827bc7b5e9bf6bba4950f8a7b37) client: document deis version
 - [`6ec3e06`](https://github.com/deis/deis/commit/6ec3e06702714f4529d4383e3ab063d062d927af) controller: add simple health check view
 - [`42464a2`](https://github.com/deis/deis/commit/42464a2d951c1890d384b0c6a2be5ed1939955da) contrib: graceful shutdown for non-ceph nodes
 - [`d7fe142`](https://github.com/deis/deis/commit/d7fe142a78b020e6b78e1397e9a2f04123a9960b) contrib/linode: allow for cluster expansion and standardize scripts
 - [`3f4e25a`](https://github.com/deis/deis/commit/3f4e25ab269eb32b7d48cea7b234a727a2555986) router: make vhost_traffic_status_zone configurable

#### Fixes

 - [`e2aeace`](https://github.com/deis/deis/commit/e2aeaced2303172cac196ceab3075d2a37c5edc5) logspout: discover logger connection continuously
 - [`3052fe6`](https://github.com/deis/deis/commit/3052fe6707708b8207cbba3e7437142545f87caf) controller: prevent overlapping config:set operations
 - [`040f90d`](https://github.com/deis/deis/commit/040f90df79648a959d1bd339163784a45f20bd89) client: only delete local ~/.deis/client.json if cancelling logged in user
 - [`67090ae`](https://github.com/deis/deis/commit/67090ae58c5cd2c3c5640162a3d2c3f542c59e10) controller: use django HttpResponse for logs
 - [`17da397`](https://github.com/deis/deis/commit/17da397b3fc9b961381ca34baed639f9e54746df) controller: use double quotes to escape ENV values
 - [`d481c4c`](https://github.com/deis/deis/commit/d481c4c35a36dd5258283636127869aef0c18a60) api: disable adding wildcard certificates
 - [`ad558e4`](https://github.com/deis/deis/commit/ad558e4bec7d3f17785e78d647267350316bba18) auth: return a 409 if a user is cancelled that has apps
 - [`2cf9205`](https://github.com/deis/deis/commit/2cf92051ee86146884c4d9bc968d3783585cde3d) setup-node.sh: update packages for a Jenkins node
 - [`46bbba7`](https://github.com/deis/deis/commit/46bbba731280cbe381b914cdd39b42b3dc3bddfb) router: Allow for comma-delimited X-Forwarded-Proto
 - [`8f0119f`](https://github.com/deis/deis/commit/8f0119f9303ad83df6a795832da163f72a89f9b0) controller: tag keys can be lowercase or capital
 - [`ba38edc`](https://github.com/deis/deis/commit/ba38edc8e4f320996780fd6c8d76336364870c55) builder: fail on piped command
 - [`78a16fb`](https://github.com/deis/deis/commit/78a16fbeb793310a55b44336df862b7662704ccd) builder: source proxy envvars
 - [`35a97b1`](https://github.com/deis/deis/commit/35a97b1a08d98b21cec4d0f022979a5db8b7e355) builder: remove temporary build dir on success
 - [`f7099c7`](https://github.com/deis/deis/commit/f7099c70a54b9dd6bf29d393e7e404d2e43c6e3b) builder: log env tcp requests information to debug
 - [`1cb11dd`](https://github.com/deis/deis/commit/1cb11ddd5d03468abefacd7bbb2494570749a934) client: simplify URL prefixing
 - [`4a9f791`](https://github.com/deis/deis/commit/4a9f791ea39f366a4941b18e0b81472d8ab2dbf3) vagrant: fix Vagrantfile to handle spaces
 - [`3275683`](https://github.com/deis/deis/commit/3275683a1a47119232f1171ded579608477fda4f) builder: demote handshake failure log to debug
 - [`e0aa2b1`](https://github.com/deis/deis/commit/e0aa2b138c558ce3372983801595f2a19e5a46b6) client: strip controller port when parsing git remotes
 - [`112f513`](https://github.com/deis/deis/commit/112f513cecf731cb22c61fe84c4c753af5429c52) controller: do not require slash at the end of the `GET /v1/users`
 - [`a450965`](https://github.com/deis/deis/commit/a4509659cb1f9352f619d2fa027e3cccb3b3a434) builder: remove empty newline
 - [`7290dd0`](https://github.com/deis/deis/commit/7290dd019d9d228a57a83f7ede812f46f13d79f3) logspout: Truncate lines too big for a single UDP packet
 - [`0700028`](https://github.com/deis/deis/commit/07000288742edd43d43a5457ca9e457341be2b80) router: Fix issues establishing real end-client IP
 - [`05e4b57`](https://github.com/deis/deis/commit/05e4b57f35cc4fb736a53f54311212474978bcb8) controller: legacy fix should modify dict instead of dict view
 - [`51c6861`](https://github.com/deis/deis/commit/51c686151ce57670f493a95da3d69e2956d96de1) database: bail out if unable to check for existing backups
 - [`8b824df`](https://github.com/deis/deis/commit/8b824df1f866491fb55dd166198ee277121b6968) client: backport deis/workflow#280 to v1
 - [`f46b967`](https://github.com/deis/deis/commit/f46b967e9fd7c2260cddad656393d7817bcd0897) tests: remove reference to python client
 - [`455e6b0`](https://github.com/deis/deis/commit/455e6b04e30c9ebcd56e25e9e2b6198e3c95efc5) tests: use and patch known "good" version of mock-s3
 - [`481fe6b`](https://github.com/deis/deis/commit/481fe6b1ce0405ef4f0eda73dcc1669edbfe2577) registry: disallow dots in s3 buckets

#### Documentation

 - [`7703d12`](https://github.com/deis/deis/commit/7703d1287d6bd879fc65465958902b6fd6ca7e56) managing_deis: recommend m3.medium for RDS
 - [`1b9facd`](https://github.com/deis/deis/commit/1b9facd36044b59abe4d9ddbd33bebb85fde220a) roadmap: add March meeting archive
 - [`48eb886`](https://github.com/deis/deis/commit/48eb88630b442f4e08c67b2c1307bd346e8b0582) contrib: add link to deis-phppgadmin
 - [`4009ffd`](https://github.com/deis/deis/commit/4009ffdec6c326c1fa9c01a5a21583548582616e) reference: remove sidenote about pushing only to master
 - [`a9a7f33`](https://github.com/deis/deis/commit/a9a7f33ac63085434e3c3aab583b2484c77a9d77) roadmap: add March release planning meeting
 - [`997e082`](https://github.com/deis/deis/commit/997e0826392ec8c03d0bcaad04c1092fc3513bc1) Add/Remove hosts: fix wrong filepath to user-data.example
 - [`ada6c9c`](https://github.com/deis/deis/commit/ada6c9c96798e96f4cba683b1282c4b6d4605a73) managing_deis: add workaround for cephless cluster
 - [`f04ed12`](https://github.com/deis/deis/commit/f04ed124f857a52c4b1e37be943e11e7839203a0) managing_deis: update Sematext agent name and URL
 - [`3faf421`](https://github.com/deis/deis/commit/3faf42182004056423039a65b3a2f8fa4470d619) installing_deis: Add parameter description about publicDomainName.
 - [`2520e54`](https://github.com/deis/deis/commit/2520e543be81f8a9f6ac1bf9f61e0eaf6d2b10f7) roadmap: add January 2016 meeting + archive
 - [`ee77077`](https://github.com/deis/deis/commit/ee77077de6ef7a979a91a63749e5659b43ef2589) roadmap: add December 2015 and January 2016 meetings

#### Maintenance

 - [`5c8d0c4`](https://github.com/deis/deis/commit/5c8d0c4bda4c3de31dd315b98561f0966115a614) reqs: update docker-py to 1.7.2
 - [`a5e065f`](https://github.com/deis/deis/commit/a5e065f6e483d7139088dc8f60c2c8f671ca5718) build.sh: remove unused git install
 - [`4ad3329`](https://github.com/deis/deis/commit/4ad3329ef6ef61f193b7da227072b7a2bbc51dbd) requirements: remove obsolete marathon lib
 - [`52e96a5`](https://github.com/deis/deis/commit/52e96a5c193677976c4809292bfa464889f55de2) (all): update base to alpine:3.3
 - [`1a79431`](https://github.com/deis/deis/commit/1a794318d0bbafa4755beb1578cd0a983774c19f) (all): bump CoreOS to 899.13.0
 - [`a7ee6cc`](https://github.com/deis/deis/commit/a7ee6cc9ad8e640ef83f63e01f46692e03ded3a4) buildpacks: update heroku-buildpack-java to v44
 - [`02fbaf1`](https://github.com/deis/deis/commit/02fbaf19df415384603b0e528e567b3c6fa7b735) buildpacks: update heroku-buildpack-php to v97
 - [`452a028`](https://github.com/deis/deis/commit/452a02824255ffd71cb41aed0898f73772b31ec6) buildpacks: update heroku-buildpack-ruby to v145
 - [`ab88aaf`](https://github.com/deis/deis/commit/ab88aafd405ac54d98de12ddd23a9526d1c90611) buildpacks: update heroku-buildpack-grails to v20
 - [`408f053`](https://github.com/deis/deis/commit/408f053ce1225ef83f8c48d9adde250be1a7bbd4) buildpacks: update heroku-buildpack-scala to v67
 - [`c23dbff`](https://github.com/deis/deis/commit/c23dbff0ed5f54c4d2710ef965d1a1e93ee59bd8) buildpacks: update heroku-buildpack-multi to v1.0.0
 - [`79d68a3`](https://github.com/deis/deis/commit/79d68a37ed5abda13bc1056099cd7eeea13503b1) buildpacks: update heroku-buildpack-nodejs to v89
 - [`4ef40e6`](https://github.com/deis/deis/commit/4ef40e6923521f313c67993e1ec96e3ae220d1c8) buildpacks: update heroku-buildpack-python to v78
 - [`ee22d78`](https://github.com/deis/deis/commit/ee22d78376f4aa90840834f5085c4ce51f2ddd59) buildpacks: update heroku-buildpack-php to v94
 - [`f448577`](https://github.com/deis/deis/commit/f44857728875cc603fa0fe820a90077671573776) buildpacks: update heroku-buildpack-scala to v66
 - [`c9b0a4c`](https://github.com/deis/deis/commit/c9b0a4c95a394808d8c46d042d06d40c65fb46bd) buildpacks: update heroku-buildpack-nodejs to v88
 - [`df95748`](https://github.com/deis/deis/commit/df957488417543043cb41bcbc051c89d05b30730) buildpacks: update heroku-buildpack-go to v31
 - [`1d8f9cc`](https://github.com/deis/deis/commit/1d8f9cc854a673e6cab1c9b29d41fb26be85ef2a) (all): bump CoreOS to 835.13.0
 - [`c96aad6`](https://github.com/deis/deis/commit/c96aad6d59157d046bfc7172716245dc9f27dd75) requirements: update docker-py to 1.7.0
 - [`96796c6`](https://github.com/deis/deis/commit/96796c64c62e5011f5352af325b0e2fb2611840c) buildpacks: update all Heroku buildpacks
 - [`3a619bc`](https://github.com/deis/deis/commit/3a619bcd974b6beb32be13643da9a47eb9d20057) (all): remove k8s scheduler code
 - [`5c916c8`](https://github.com/deis/deis/commit/5c916c8bd9f75669234c851284c5505cef7691e6) Godeps: bump googleapi, remove unused packages
 - [`7c37363`](https://github.com/deis/deis/commit/7c373633f1aaf05741cb5c8f745de587db17212c) (all): bump to CoreOS 835.11.0
 - [`4dac497`](https://github.com/deis/deis/commit/4dac49741901f1a780693d464eea1981ce29f68b) deisctl: update stateless warning message
 - [`d88ec07`](https://github.com/deis/deis/commit/d88ec07c99f4ab903fc7d250025146b966aacee8) (all): bump to CoreOS 835.9.0
 - [`31ad174`](https://github.com/deis/deis/commit/31ad17416487d394a8d06f6b262d9ed4eb807020) controller: update docker-py to 1.6.0
 - [`aad62fe`](https://github.com/deis/deis/commit/aad62fedd20acc567d1477c152884de77397ffa1) (all): bump CoreOS to 835.8.0

### v1.12.2 -> v1.12.3

#### Maintenance

 - [`7356d26`](https://github.com/deis/deis/commit/7356d26adeb2f48d8df82a1d0baa090da7d4bb20) (all): bump CoreOS to 835.12.0

### v1.12.1 -> v1.12.2

#### Fixes

 - [`0738c13`](https://github.com/deis/deis/commit/0738c13949187a6d444a693193585aed7e44304a) database: supports HTTPs as S3 endpoint
 - [`3498099`](https://github.com/deis/deis/commit/34980992c27497735b9155ce7ff8786c048e1f41) router: fix router common prefix app publishing
 - [`27eab71`](https://github.com/deis/deis/commit/27eab713f9dd1c2a4f3bc5b8152b43df7bec66b5) contrib: Add drop-in to make docker require flannel
 - [`42d00af`](https://github.com/deis/deis/commit/42d00af18390745da4bebf6b3f70c11d40e09884) deisctl: don't panic when config key/value is malformed
 - [`7410fb7`](https://github.com/deis/deis/commit/7410fb74c7fe5c85a9df1765d5b18ea836e68f8d) builder: Fix problem with missed git repos after builder restart

#### Documentation

 - [`a19caaf`](https://github.com/deis/deis/commit/a19caaf61ad7e3f718570b93f9fe11aa0665fa2f) managing_deis: change swift3 link.
 - [`82732b4`](https://github.com/deis/deis/commit/82732b486f4531242b41b44493c7f69da3aa2425) hacking: add docs to show how to use the docker-machine env

#### Maintenance

 - [`0f96abe`](https://github.com/deis/deis/commit/0f96abea0a7d28257f8c87e89eb6453fec21fa79) contrib/coreos: remove debug-etcd service
 - [`e1e3927`](https://github.com/deis/deis/commit/e1e39274cbed8e5fc04dfa5b6a42dc74ae1b61f7) MAINTAINERS: don't enumerate maintainers

### v1.12.0 -> v1.12.1

#### Fixes

 - [`293d657`](https://github.com/deis/deis/commit/293d6572d2d0da73595bd7e178e438e83115c2e6) registry: fix create_bucket s3 compatability
 - [`b07db69`](https://github.com/deis/deis/commit/b07db69e9b3ad197d869fec4cb0b5a7051a0ef16) user-data: always start flannel on boot
 - [`d9ef023`](https://github.com/deis/deis/commit/d9ef0234f4cbab8f806f7970d71d883da1261e52) contrib: re-introduce data dir mapping for etcd
 - [`070d081`](https://github.com/deis/deis/commit/070d081838e262256b7987fd2e9e6c44aff48dd1) create_bucket: check for existence of None

#### Maintenance

 - [`3eb277a`](https://github.com/deis/deis/commit/3eb277a4f90c2bd2733764e49b9ce70e7c3567e9) (all): bump CoreOS to 766.5.0

### v1.11.2 -> v1.12.0

#### Features

 - [`c8100da`](https://github.com/deis/deis/commit/c8100da71e24a2b932b453c886a8cfaa904327b5) contrib/digitalocean: improve Digital Ocean provisioning workflow
 - [`5c8f5b9`](https://github.com/deis/deis/commit/5c8f5b984b361e7379d7bf5048ae46abcfa0e800) tests: use Docker 1.8.3 for Jenkins nodes
 - [`a47fcfe`](https://github.com/deis/deis/commit/a47fcfe4586e43e71c4515075626c3e0814a07d1) contrib: vagrantfile sources utils.sh
 - [`25653e3`](https://github.com/deis/deis/commit/25653e39a1ac8d76eee28563e2601c9750c3e8a7) builder: make SOURCE_VERSION env variable
 - [`d55f9a1`](https://github.com/deis/deis/commit/d55f9a14bef9ef72b7c365f04a3461d4c6c816c8) registry: add more information on during startup process
 - [`8f7f284`](https://github.com/deis/deis/commit/8f7f284895601087842e513855cc62072bda627c) controller: disable swap if there's a mem limit
 - [`ab53096`](https://github.com/deis/deis/commit/ab53096daf466024cda23c203b377326815df4b9) controller: use Docker for last-mile layer
 - [`9e87112`](https://github.com/deis/deis/commit/9e87112957f5fdaf852028df018548e2641937d4) deisctl: bump fleet api to v0.10.2

#### Fixes

 - [`2ecc28f`](https://github.com/deis/deis/commit/2ecc28fe8e0e791c8903aea99a5897a54f4543e0) database: revert to alpine:3.1 for PostgreSQL
 - [`2a0a319`](https://github.com/deis/deis/commit/2a0a319cf5d3d77f95ba92f7d34062390b41bb72) contrib/linode: fix discovery url issues in deployment scripts
 - [`61c1026`](https://github.com/deis/deis/commit/61c1026f480e73c11a06bebccdee30bc44ed47a4) builder: continue updating etcd after errors
 - [`055763b`](https://github.com/deis/deis/commit/055763b9d2e03d7fc5c87ac46f2bc3d5a013fc24) client: init procfileMap in "deis pull"
 - [`a81e689`](https://github.com/deis/deis/commit/a81e6896137e4775e6799d3b99f9c0a278488112) contrib/linode: fix issues with DHCP on Linode
 - [`2c5d0da`](https://github.com/deis/deis/commit/2c5d0dad2e10508d780a5721d0f9da371804c7ad) Vagrantfile: specifically use bash to source utils.sh
 - [`8ac36e8`](https://github.com/deis/deis/commit/8ac36e8bf5107e1f9efbb05f138d7b631b1f78e6) client: catch and propagate client errors
 - [`5fc33de`](https://github.com/deis/deis/commit/5fc33de0d3e291c9b2a8d52241d6e52f28d11196) contrib/gce: replace discovery url properly
 - [`d830adf`](https://github.com/deis/deis/commit/d830adf1e8d323ff37f23300661c543daebc632b) contrib/azure: use Azure premium storage
 - [`648439d`](https://github.com/deis/deis/commit/648439dcdcc0a9b8ce7de268fdcf394cc2f1e758) contrib/azure: discovery url
 - [`4d5db42`](https://github.com/deis/deis/commit/4d5db42da2b27e825af2d3aa7548dda6dc3957a9) builder: update d-in-d wrapper for newer Docker
 - [`b9dde4a`](https://github.com/deis/deis/commit/b9dde4a90fdc6985479b6b2f9cd1dcad5195326d) client: read procfile from PWD when pulling
 - [`87afcd9`](https://github.com/deis/deis/commit/87afcd91e8b98c1dc9bd351496aaab73da4376da) controller/web: enable unit tests
 - [`df1dd36`](https://github.com/deis/deis/commit/df1dd3692b97f9b8463b7aebd9baf409fd8e9fdd) contrib/aws[rigger]: upload correct ssh key

#### Documentation

 - [`70bd364`](https://github.com/deis/deis/commit/70bd36465bb0529480324c30071611ee411c9d9b) roadmap: add release criteria
 - [`3ec1ccf`](https://github.com/deis/deis/commit/3ec1ccfaf19e3eb9fac06815a59ec01c378a7201) contributing: update where build instructions are
 - [`393883e`](https://github.com/deis/deis/commit/393883e2d35b185164820d35e3ac62a4feba8191) contributing: add optional pre-reqs
 - [`dfa98e9`](https://github.com/deis/deis/commit/dfa98e9df93ce249af3320f1f1015b49c881d5db) upgrading-deis: graceful upgrade in stateless mode
 - [`c489394`](https://github.com/deis/deis/commit/c489394f0dd49e810f74c1513aa6369bb6fa380e) roadmap: add October 2015 planning meeting

#### Maintenance

 - [`386947f`](https://github.com/deis/deis/commit/386947f14d96025d12eee3c09980125a90021482) router: update nginx to 1.9.6
 - [`ea93351`](https://github.com/deis/deis/commit/ea9335140d906d11f6e5e1496d042d2df2240e9a) logger: improve drains
 - [`491bd26`](https://github.com/deis/deis/commit/491bd269f9e398402c5bef30c5b2b76a70763b54) builder: update Docker to 1.8.3
 - [`f78a706`](https://github.com/deis/deis/commit/f78a706bca18c4f334470a1ae799434633bb20ec) (all): update Docker base image to alpine:3.2

### v1.11.1 -> v1.11.2

#### Fixes

 - [`abe3fa0`](https://github.com/deis/deis/commit/abe3fa09ebda348f8c2c91113571bc026bc2eaf6) fix(user-data): use $private_ipv4 for fleet
 - [`c55d9c8`](https://github.com/deis/deis/commit/c55d9c8bdab91d1539fa34fbb7ec392c7355ce24) fix(*) add flocks to additional services
 - [`6dc0da4`](https://github.com/deis/deis/commit/6dc0da42e096ea6132e2bc8ffc861cde8e1bb937) fix(*) add flocks to -data pulls for db & builder

### v1.11.0 -> v1.11.1

#### Fixes

 - [`7c02cd6`](https://github.com/deis/deis/commit/7c02cd6bfa51a3a6720bfa3d03a403ea7d16a117) fix(*): use flock to serialize certain pulls

### v1.10.1 -> v1.11.0

#### Features

 - [`46a5ec7`](https://github.com/deis/deis/commit/46a5ec7df73a8015a1687ee0c36a490d3cf79003) Registry: add support for IAM role authentication
 - [`b60b8f9`](https://github.com/deis/deis/commit/b60b8f9932ba3be26ea2b643c36c37e1ba862a68) client: print error messages to stderr
 - [`4d25e8b`](https://github.com/deis/deis/commit/4d25e8b300024b27d861f967580fc8c4f122115e) contrib/vagrant: auto discover DEISCTL_TUNNEL via Vagrant
 - [`ff221e4`](https://github.com/deis/deis/commit/ff221e4954871b3da590d1c97f04c3e3c154a32a) controller: publish etcd keys at controller start
 - [`4a5606f`](https://github.com/deis/deis/commit/4a5606f50c6f6f26e5b716ec51bcb687bc57f78c) router: improve healthcheck
 - [`3838a0a`](https://github.com/deis/deis/commit/3838a0ad238baba12bc3128a8c8073431563cbc8) Router: IP whitelisting for apps and controller
 - [`418b998`](https://github.com/deis/deis/commit/418b9986ca2003982db31cae772978c2c3807734) controller: allow users to transfer app ownership.

#### Fixes

 - [`e19fdcb`](https://github.com/deis/deis/commit/e19fdcb580c2fef2cf2f07211f61dd4c192e9be7) client: add newline for usage message
 - [`b11b45c`](https://github.com/deis/deis/commit/b11b45c53091e6bd9926b213a61b77e19017a9d9) deisctl: remove unused options docstring
 - [`29b6c0b`](https://github.com/deis/deis/commit/29b6c0b2aab51739cd2a128f7fd5dbbb6b305809) contrib/aws: stray bracket in aws query removed
 - [`6e91f50`](https://github.com/deis/deis/commit/6e91f5034c3aa7bfff6110f2418b7c85134f7ff0) tests: binaries are installed in the wrong directories
 - [`7421f01`](https://github.com/deis/deis/commit/7421f01cbcd0c6be617b9c4c8f9e4408ef63f8b0) client: accept multi line config vars.
 - [`1b331a9`](https://github.com/deis/deis/commit/1b331a9172f999f0795e56f9b97e26457b5bfd16) deisctl: fixes for graceful upgrades.
 - [`3a287ac`](https://github.com/deis/deis/commit/3a287ac722c060284f8e63c32275b71f770de566) contrib: Zone argument not supported for aws
 - [`01daf4d`](https://github.com/deis/deis/commit/01daf4dbbfec07f2ce73b3cbfd3106107db29364) contrib: urlopen call hangs
 - [`799e559`](https://github.com/deis/deis/commit/799e55995ed630a93526c17a67dc81932b9b0d8d) deisctl: fix refresh-units as cache is now missing
 - [`306e63d`](https://github.com/deis/deis/commit/306e63d4c9aeaaa793727af8faa4ec70b1811db2) controller: validate config keys at the api level
 - [`e3a28b4`](https://github.com/deis/deis/commit/e3a28b48a451338e12165ab769d9c711f9533805) client: suggest help when git remote creation fails
 - [`0bd2471`](https://github.com/deis/deis/commit/0bd24715f24d23b705ec4cd0a8d266506b95c4f5) contrib/coreos: mask etcd.service
 - [`f7d55f6`](https://github.com/deis/deis/commit/f7d55f66bdc178648d6f5abaabbbaf4025d165c8) deisctl: add status command to help screen
 - [`53bf1ed`](https://github.com/deis/deis/commit/53bf1ed8ec2b44f34290aa04318734a779214960) controller: state property can't be used in a list_filter in the Django web admin
 - [`05b6d03`](https://github.com/deis/deis/commit/05b6d03dc7f01251eb82fe96d2814de38d88c3f2) client: print logs that lack a category prefix
 - [`b12772d`](https://github.com/deis/deis/commit/b12772d3bd09503c345c5a7403b58582ffb039b8) firewall: check etcd is running
 - [`404889c`](https://github.com/deis/deis/commit/404889c9a835d442818592d83cf0f2a2d22f2fc4) controller: write app event logs via logspout
 - [`c5bb64b`](https://github.com/deis/deis/commit/c5bb64ba5851cdc4ae97c43de1523ad43f509499) Makefile: silence docker-machine warning when it isn't installed
 - [`7bd72ff`](https://github.com/deis/deis/commit/7bd72ff63daf277962fbcf4edfd0d6770d5db37b) controller: healthcheck no longer check non-web processes
 - [`e68d5e2`](https://github.com/deis/deis/commit/e68d5e292814231e30054d98dab2df47e4278928) client-go: error deserializing json
 - [`1330404`](https://github.com/deis/deis/commit/1330404b3f84206c71294ac2d65c050e01b043b0) controller: clarify error if registration is disabled

#### Documentation

 - [`c3381a0`](https://github.com/deis/deis/commit/c3381a0cca1add9fd62535e74f8f885911d75c51) installing_deis: remove try.deis.com references
 - [`70b060c`](https://github.com/deis/deis/commit/70b060c001ed952ac2371266d89362a0c8971265) installing_deis: install Azure CLI
 - [`bc33750`](https://github.com/deis/deis/commit/bc337508b025595e135b087052bbedd40ceafb46) managing_deis: add caveat for different etcd versions
 - [`c78854f`](https://github.com/deis/deis/commit/c78854fcd458357af58593030943b997616d8080) managing_deis: add OpenStack Swift store example
 - [`751f6cf`](https://github.com/deis/deis/commit/751f6cfb7d237e5f88b23b069dd0ba480084f5b4) understanding_deis: add section on router mesh
 - [`55cb83f`](https://github.com/deis/deis/commit/55cb83ff854eae1257145fd8cc6820456cb6d224) running-deis-without-ceph: s3 region var
 - [`f9cd978`](https://github.com/deis/deis/commit/f9cd978ae1ab7c696d2b74cde211b6c3cd2afb3f) quickstart: added install deisctl step at top of quick start docs

#### Maintenance

 - [`2333ccd`](https://github.com/deis/deis/commit/2333ccdc19969a6affa1b8b2225b88d186a2af33) logger: refactor-- also includes new features
 - [`6026ac6`](https://github.com/deis/deis/commit/6026ac61c12070f098e21ba978a22b970e7f5a15) tests: update shellcheck to 0.4.1
 - [`c2dedfb`](https://github.com/deis/deis/commit/c2dedfb13546c1840a043a49f50fc657f3df1adb) (all): update CoreOS to 766.4.0
 - [`43f0833`](https://github.com/deis/deis/commit/43f08330bb4fcf830447a8d962b0f25238c250e1) (all): update CoreOS to 766.3.0
 - [`a1b431b`](https://github.com/deis/deis/commit/a1b431bcfa57e5768f980110e3ffa66b09507784) tests: update virtualbox and vagrant

### v1.10.0 -> v1.10.1

#### Fixes

 - [`04a9306`](https://github.com/deis/deis/commit/04a93066c6dc7bd702e72d6bf23870ceed3b7446) deisctl: default stdout and stderr are reversed
 - [`5878b8f`](https://github.com/deis/deis/commit/5878b8f8c3d7d625351d93589cf8223d06a93bbb) (all): ensure component containers stop
 - [`9ef640f`](https://github.com/deis/deis/commit/9ef640f7ca06b9c146d3525d9f6bb5a09f9db40d) database: improve postgres ready checks
 - [`fb44cbc`](https://github.com/deis/deis/commit/fb44cbc7d4de877f1487498fbfe56f3e284b7db8) database: fix db init file perms

### v1.9.1 -> v1.10.0

#### Features

 - [`540b9f7`](https://github.com/deis/deis/commit/540b9f7f9b9ad4c68b3f749f8ffc0924142bbe77) Makefile: retry discovery-url five times
 - [`cdd6911`](https://github.com/deis/deis/commit/cdd691133287430955674e35a29e9712ab09ca62) client: add makeself installer
 - [`30fefce`](https://github.com/deis/deis/commit/30fefce3c5b3cebad6bb58aa1c6b4a12eeb0e822) logspout: support sending log via tcp
 - [`ab4f091`](https://github.com/deis/deis/commit/ab4f091d7b6638b202fd6e078d4f7de05d536429) builder: remove Dockerfile shim
 - [`5c5f822`](https://github.com/deis/deis/commit/5c5f822000d757a1f6e2569ae3692dc4b11f4a3e) builder: reimplement the builder in Go
 - [`d3a333f`](https://github.com/deis/deis/commit/d3a333f95ac39a852f71d8e4b2ac9e0eb5894b95) router: improve nginx status page
 - [`bbf5766`](https://github.com/deis/deis/commit/bbf57666d7129e354c6bdb6cf1e6598ca8a9241f) Makefile: option to disable store components
 - [`fbe9b84`](https://github.com/deis/deis/commit/fbe9b8477c978c992c9de100391683705a3fb172) client-go: support deis plugins
 - [`8d0637c`](https://github.com/deis/deis/commit/8d0637c3430313b3188d548c2f7291caf745ac36) controller: allow admins to delete users
 - [`861052f`](https://github.com/deis/deis/commit/861052fc33d6ba431ee1154d885bff9caf7d5dc8) controller: allow admins to change a user's password.
 - [`a5dc3b4`](https://github.com/deis/deis/commit/a5dc3b4ddf5c370ddee96fa0c392dcd822e33549) platform: support placement options for each plane and router mesh
 - [`3896107`](https://github.com/deis/deis/commit/3896107695e322b24b774105f5e2c5316b8e3239) client-go: add minimal windows support
 - [`313ab04`](https://github.com/deis/deis/commit/313ab0452103ed76e22e7ec9f3b8dc33c18be577) (all): replace /bin/bash with /usr/bin/env bash
 - [`fa57b1f`](https://github.com/deis/deis/commit/fa57b1f10e1169d1a8ea2b9673c81d81d3664145) router: enable ssl caching and expose ssl_protocols for configuration

#### Fixes

 - [`8232455`](https://github.com/deis/deis/commit/823245526d223a1855fba6b9e05f9fa0f8a0a7bf) database: arrange etcdctl arguments properly
 - [`88e6f7a`](https://github.com/deis/deis/commit/88e6f7a2ca6c351406d2b7bd4ff7871c0ba3b9f2) database: use ETCD host during database recovery detection
 - [`6839b59`](https://github.com/deis/deis/commit/6839b592a1ff9c468a8c91f5f2b6817aebd5ae86) logspout: set GOMAXPROCS=1
 - [`5c36bca`](https://github.com/deis/deis/commit/5c36bcab5410b4d1cdab75b11e20db919bd5b9ef) client: show non-JSON HTTP errors
 - [`3b2827a`](https://github.com/deis/deis/commit/3b2827a4ab140023d2caf68b663f7be1c5bc81e0) deisctl: error when journaling a global unit
 - [`1f54bb5`](https://github.com/deis/deis/commit/1f54bb55c4904b64655c7f47e9460e8fb5dfec09) client: fix make install
 - [`bd3390f`](https://github.com/deis/deis/commit/bd3390f3958f7168d481677fa16b970ee906e8bd) client: print usage if no args were given
 - [`78ab6cf`](https://github.com/deis/deis/commit/78ab6cfc7b3895f0c264171fd1d71740fc609b77) router: revert symlink to /dev/stdout
 - [`6bcd7d7`](https://github.com/deis/deis/commit/6bcd7d78e90c8ea21e6ef9dc1611e9e49686526c) client-go: pass arguments when using root command to list
 - [`ea94180`](https://github.com/deis/deis/commit/ea94180c7fff82cf8eb6dd39bc263060bfdcd0c1) logger: add unique id for each app name in logs
 - [`c082105`](https://github.com/deis/deis/commit/c082105c666537ab93422e2d37bc7a11eb8e660c) builder: added back the check-repos script
 - [`eb3f7e6`](https://github.com/deis/deis/commit/eb3f7e6b26970dbc969de8f75044e77af943b86b) docs: update NewRelic section.
 - [`4efdc61`](https://github.com/deis/deis/commit/4efdc616f642e9dab372009d3fb8d25a6bcf6d60) deisctl: avoid stdout issues in tests
 - [`4467aad`](https://github.com/deis/deis/commit/4467aad156f0da8571629e617de501d23c49c655) logger: allow underscores in container names
 - [`960b4fe`](https://github.com/deis/deis/commit/960b4fefed8c09e1db4351c6c78377031953a7a0) controller: don't check permissions with @permissions_classes
 - [`bd5f008`](https://github.com/deis/deis/commit/bd5f008bb53ae255c1eda1520cb68202a2ff70fd) logspout: fix deadlock when attaching
 - [`ae64f49`](https://github.com/deis/deis/commit/ae64f49045352147f296948d534d0a9d4df76942) contrib: AWS instances may only have a private IP
 - [`7f72c70`](https://github.com/deis/deis/commit/7f72c70a1c5b6011c09723f7c0ac0bd6b7512105) client-go: paginate over lists
 - [`bfb8844`](https://github.com/deis/deis/commit/bfb884475109f3f409beae28e5dfcae3ef575241) deisctl: clarify registration message
 - [`d6289ba`](https://github.com/deis/deis/commit/d6289baebf9bb65edfe8ff77d5643c5d620848a0) contib: remove nonfunctional scripts of dubious usefulness
 - [`e2c0689`](https://github.com/deis/deis/commit/e2c06895ec5202a751cfb2b61f8b64039de98a2e) contrib: remove manage_etc_hosts from example user-data
 - [`7700196`](https://github.com/deis/deis/commit/77001967a38f84f25a3b32ec3dd09df93155d6f4) deisctl: improve dock name match
 - [`d83a3e8`](https://github.com/deis/deis/commit/d83a3e8bb8abd7336567cb04def906330c88d0d5) deisctl: fix deisctl ssh <cmd> output
 - [`e5565a2`](https://github.com/deis/deis/commit/e5565a2bd424962062450db9fcc7e45af6649172) client-go: unset token before logging in
 - [`e73447a`](https://github.com/deis/deis/commit/e73447a782eed0cd83a20b7c68a96d8485c72c9e) client-go: fall back on directory name if remote can't be detected
 - [`4c37adc`](https://github.com/deis/deis/commit/4c37adc1d821af9298d1afcdd3af820707cb8ff3) contrib/aws: clarify proxy protocol error
 - [`381e7ca`](https://github.com/deis/deis/commit/381e7ca4143ba7593e1b0fd12f48995206b0642b) client-go: fix webbrowser merge error
 - [`05afd5c`](https://github.com/deis/deis/commit/05afd5c1a19398ba50603c705f04fad0675cfedd) deisctl: fix how target names resolve
 - [`2f5171f`](https://github.com/deis/deis/commit/2f5171f185cb44ec496c3cffb723ba083b4e22f5) Makefile: don't build swarm by default
 - [`bbcb957`](https://github.com/deis/deis/commit/bbcb95785de9f1f52a880b1ba2c276a5f7d07805) deisctl: collapse component groups
 - [`376a098`](https://github.com/deis/deis/commit/376a098ccb9ddad9912a2b143ca94dc4ae9deae4) tests: use command line flags to test auth:passwd
 - [`6f765d4`](https://github.com/deis/deis/commit/6f765d42fecc8f85041c262b418b0215c0cc0575) tests: abort if warning is found in output
 - [`941c481`](https://github.com/deis/deis/commit/941c481ad560a90aaf10ea18e7968e67bcab076b) builder: always configure SSH key

#### Documentation

 - [`97a6c2b`](https://github.com/deis/deis/commit/97a6c2bf7ec416564bf357ca43a14ed4d7e6f46e) roadmap: add September 2015 release planning
 - [`20b64c2`](https://github.com/deis/deis/commit/20b64c206d41099d60373ff462d866313b338e33) roadmap: add September planning meeting
 - [`52733f4`](https://github.com/deis/deis/commit/52733f4a13a8be3402af6ed77270106fe561dfb4) MAINTAINERS: add Seth as a contributing maintainer
 - [`b25f24b`](https://github.com/deis/deis/commit/b25f24bab2cfbbc812c261f17e47d63ac65cd6cb) azure: add vnet creation instructions
 - [`7d3db03`](https://github.com/deis/deis/commit/7d3db0316ef037d71f0ddfb749b8fdc10c652403) hacking: offer DigitalOcean credits to new contributors!
 - [`a08d433`](https://github.com/deis/deis/commit/a08d4333f306ad8f41f5061eb45e1c34bc3e8b05) understanding_deis: add documentation on node failover
 - [`8ccb760`](https://github.com/deis/deis/commit/8ccb76044f06f6a0ee3566868334009198375367) isolating-etcd: add link to user-data.example
 - [`9b45fe1`](https://github.com/deis/deis/commit/9b45fe1f8e234a58b0b963244ff58ff3adf3cda7) installing_deis: clean up old references
 - [`0246ada`](https://github.com/deis/deis/commit/0246ada1fc813dc06da014988b3ff7111990f376) installing_deis: specify CoreOS for bare metal
 - [`034709b`](https://github.com/deis/deis/commit/034709b99228083096d1c5b3599a50ed785cd532) managing_deis: remove a space from a bash substitution
 - [`3845216`](https://github.com/deis/deis/commit/3845216b94b40197ff4aa7e04cb795db6dd98f6c) client-go: update readme
 - [`b90c2fb`](https://github.com/deis/deis/commit/b90c2fbb5e4cf9a93c649b7285c85eb301a639c8) MAINTAINERS: add Wael as a contributing maintainer
 - [`eeed28a`](https://github.com/deis/deis/commit/eeed28a22cc702d73b0d365decb5288390aaf163) aws: require deisctl before provisioning
 - [`32076e2`](https://github.com/deis/deis/commit/32076e274d12fdfa5ce57018b33dd6ad8d273453) deisctl: correct isolation doc mistake
 - [`58f0748`](https://github.com/deis/deis/commit/58f074894f6afafa07c884f149585e93d1b3f3c2) roadmap: add August 2015 planning meeting
 - [`620e548`](https://github.com/deis/deis/commit/620e5482fce3921e6ca15e7a8233db91e6629897) contributing: sync up maintainer lgtm policy

#### Maintenance

 - [`56256c5`](https://github.com/deis/deis/commit/56256c5adf1b9cb46aec5b6dea733ce7758a1fc2) swarm: update to go 1.5
 - [`7b63dd3`](https://github.com/deis/deis/commit/7b63dd3c871fe118b91947e897b5bddd08b7021a) tests: use go 1.5 toolchain
 - [`101bbba`](https://github.com/deis/deis/commit/101bbba3788043e64545a2c67ead224c4958b048) router: update nginx to 1.9.4
 - [`c6d2c4a`](https://github.com/deis/deis/commit/c6d2c4a123d0311bb8eaee415b13a0029273b116) docs: Fixed rendering of a few code snippets
 - [`dd4d777`](https://github.com/deis/deis/commit/dd4d777072ec172d9b3dbca7b59da5759b422f3b) docs: use Docker Toolbox, not boot2docker
 - [`bb09e27`](https://github.com/deis/deis/commit/bb09e276ca6f3c4b07442f3bd328f6cd58e0a360) router: update nginx to 1.9.3

### v1.9.0 -> v1.9.1

#### Fixes

 - [`d9a09b0`](https://github.com/deis/deis/commit/d9a09b003112a0f4a6f882ff18413332e179ba7c) contrib: fix etcd2 data directory on AWS
 - [`8f50e70`](https://github.com/deis/deis/commit/8f50e7063a4d8669b77592edfef71f910de6e8a3) publisher: ignore healthcheck values if unset
 - [`70867aa`](https://github.com/deis/deis/commit/70867aa69e0f02fd416cd250e526e57a2a01b9a1) controller: require fleet.socket
 - [`1f2264c`](https://github.com/deis/deis/commit/1f2264cf303790bba6f66181dfcc209ab090e392) flannel: use default iface for starting flannel except vagrant
 - [`ee15c14`](https://github.com/deis/deis/commit/ee15c149270669be2eaae677f76dda7fa4874a3a) deisctl/units: stop k8s services without errors
 - [`e92e2db`](https://github.com/deis/deis/commit/e92e2db5ddd4e740b4e43d8cf0e14680be5e8bcf) contrib: fix debug-etcd
 - [`837ef9d`](https://github.com/deis/deis/commit/837ef9dc7882e5264ea343f2548afe2d6325171d) mesos-marathon: change instances to zero instead of scale to zero

#### Maintenance

 - [`76571aa`](https://github.com/deis/deis/commit/76571aa87ae423537ca240a86d5f128b026463a2) (all): bump etcd to 2.1.2

### v1.8.0 -> v1.9.0

#### Features

 - [`02830a6`](https://github.com/deis/deis/commit/02830a6684c7ffcd10d19107c52963cd51075eb4) deisctl: add graceful upgrade
 - [`cd794aa`](https://github.com/deis/deis/commit/cd794aaa8866d4da8df5c0847530279ee402fccc) platform: add http liveness and readiness health checks
 - [`15330f1`](https://github.com/deis/deis/commit/15330f14de3859526992f90285c73e401d29751d) client-go: decode json error messages
 - [`7e3c374`](https://github.com/deis/deis/commit/7e3c3748ba84f776c976b8a0cbbc38c7386b8d23) k8s:  add k8s scheduler support to deis
 - [`176579d`](https://github.com/deis/deis/commit/176579dc86a0c8e9d797903bd3141a9e557b3dc6) commit-msg: don't prevent a commit, only warn
 - [`11cdd5c`](https://github.com/deis/deis/commit/11cdd5c2d93a5348d903f645585719595906c8c3) deisctl: extend SSH to allow exec and docker
 - [`e998b29`](https://github.com/deis/deis/commit/e998b29831a34027313488e4b0856f4abe5a5822) deisctl: change deis-mesos-zk@1 to deis-zookeeper
 - [`b7f4032`](https://github.com/deis/deis/commit/b7f40322d3ddf5407ad4f2d62885543220a3b8a9) client-go: add certs endpoint
 - [`8135ead`](https://github.com/deis/deis/commit/8135ead2b4d83c85e2a539325568b8deedc750a7) deisctl: specify # of routers to install
 - [`15f1224`](https://github.com/deis/deis/commit/15f12244957a8672888f309e5e301f26a5b03906) platform: add overlay network
 - [`ff73695`](https://github.com/deis/deis/commit/ff73695f0b33ee503c11fb700f81de896b81a320) client-go: add perms endpoint
 - [`d5e56da`](https://github.com/deis/deis/commit/d5e56da3b483f342cf4f5b5e30318362e442e685) router: don't copy Godeps, godeps can find it automatically.
 - [`d82b975`](https://github.com/deis/deis/commit/d82b97564e2651638349289b7070b52dfb718d9a) contrib/ec2: allow easier selection of CoreOS version
 - [`31578ed`](https://github.com/deis/deis/commit/31578ed281844d3b7e14b7150685d39c2337db50) client-go: add releases endpoint
 - [`6eb1e43`](https://github.com/deis/deis/commit/6eb1e43eef0c808c1e350ed5e531d8f066e23631) client-go: add limits endpoint
 - [`73246fc`](https://github.com/deis/deis/commit/73246fc62b2cb3e76f3908b9a4014b19110af4e7) database: backup the db on unit stop
 - [`56b5860`](https://github.com/deis/deis/commit/56b586007de3af507e758f88523338413cf09eb0) client-go: add git endpoint
 - [`c1b4b30`](https://github.com/deis/deis/commit/c1b4b30d3a99b04b52f1c8d39549e6e8a5ca06db) client-go: add tags endpoint
 - [`bd230bb`](https://github.com/deis/deis/commit/bd230bbe7e3cc6666d7fb61f4f683e5bd3169eaf) mesos-marathon: add marathon framework support for deis
 - [`ba24426`](https://github.com/deis/deis/commit/ba24426b4ba7bb8ab0187ede194bea28318198d1) client-go: add config endpoint
 - [`b28e2df`](https://github.com/deis/deis/commit/b28e2dfdc8d44db42ea17c16f688f2e1877b1a21) client-go: add builds endpoint
 - [`2c6dfe4`](https://github.com/deis/deis/commit/2c6dfe47b2ed9ecccfbc42f630ed66952a72c900) client-go: added ps endpoint
 - [`629b8bc`](https://github.com/deis/deis/commit/629b8bcf0901ce3bf46268763d68065c72359e04) client-go: add domains endpoint
 - [`94a9fcf`](https://github.com/deis/deis/commit/94a9fcf08de0d2f0431a6c406314b24288f87902) client-go: don't print static binary check
 - [`6bc5a0b`](https://github.com/deis/deis/commit/6bc5a0bdb3ddffe2200242b8fbc94fe0987c6207) pkg/prettyprint: support colorizing terminal text

#### Fixes

 - [`1b9d1b9`](https://github.com/deis/deis/commit/1b9d1b99fc9389137f67d3f8c55f3e8e0a22fe4d) swarm: fix swarm shutdown
 - [`510eb7d`](https://github.com/deis/deis/commit/510eb7d464ed6c599084bde35ff6b4ff34f83b43) mesos: quote properly in marathon build
 - [`9cba688`](https://github.com/deis/deis/commit/9cba688aa33bea301fdc32526dd885ece3cfa03a) contrib/openstack/provision-openstack-cluster.sh: fix openstack …
 - [`518044a`](https://github.com/deis/deis/commit/518044ac97e620955f887496bac559c939658d21) contrib/linode: Install CoreOS from the stable channel
 - [`c98183a`](https://github.com/deis/deis/commit/c98183a9476ae5b7404586579f6a84293ea34194) k8s: deploy app into individual namespace
 - [`bf94812`](https://github.com/deis/deis/commit/bf94812f59f171a072bb4f2952342de340172af1) client-go: accept all non error HTTP codes
 - [`3a94fec`](https://github.com/deis/deis/commit/3a94fec5404394eecfbd6a64d83fca001b1d4b60) controller: prevent < 1 gunicorn workers
 - [`d568115`](https://github.com/deis/deis/commit/d568115fc363ba12bbfd514eaf18aec1d889b111) deisctl: allow help for ssh and dock
 - [`3c457bf`](https://github.com/deis/deis/commit/3c457bf9ca232f916df40bb4817396b8b20b37e8) database: fix if statement rather than ignore shellcheck error
 - [`9054447`](https://github.com/deis/deis/commit/9054447422627e434a8dc879042f59760f177fe9) deisctl: add dock to the list of commands
 - [`f921345`](https://github.com/deis/deis/commit/f9213452634072bd43fe5dab1c2a6e327ddce311) builder: fix spelling mistakes
 - [`aaeb33e`](https://github.com/deis/deis/commit/aaeb33e0cc6b63e0605ddeaf1c33e84ffc31f7a0) controller: remove extra bracketing
 - [`ae88f2a`](https://github.com/deis/deis/commit/ae88f2af5e33838edd68269220de95f7b6573bc5) k8s: change etcd2.service in deis-kube-apiserver.service
 - [`f886695`](https://github.com/deis/deis/commit/f886695bd98d3063658adf770f2dedfbd842e4cc) tests: fixup deis binary lookup
 - [`80ff17f`](https://github.com/deis/deis/commit/80ff17ffdaf4df895b4ac5ac0e8789f788bdd21f) client: remove nonexistant certs:update
 - [`55b9f0f`](https://github.com/deis/deis/commit/55b9f0f67873cf987dd8c89413ebf98baf7f275c) client-go: make output from apps:destroy match python client
 - [`de6e714`](https://github.com/deis/deis/commit/de6e71464f5afe7ce3f07c70351be2f627a8c3d7) vagrant: fix too many redirects problem
 - [`8a672ad`](https://github.com/deis/deis/commit/8a672ad71c62919f4ebdb70e2cee874d25d57db4) deisctl: exit when stop/start fails
 - [`80c2f08`](https://github.com/deis/deis/commit/80c2f0847908528dc004941394b3b3d7303a14d8) controller: convert config to correct types
 - [`4336ef2`](https://github.com/deis/deis/commit/4336ef2adebde3e8b38b7ca7daf3d25ef68ae92f) mesos: build images with prefix and tag
 - [`f6d39c0`](https://github.com/deis/deis/commit/f6d39c0b6ce32aab02818e7be533e4a3630bfde8) builder: until push-images script exists
 - [`0bd2568`](https://github.com/deis/deis/commit/0bd256829cc307feef75d6166c2c36d944e63851) Dockerfiles: silence apt-get TTY warnings
 - [`a045639`](https://github.com/deis/deis/commit/a0456397cc39fc4a2769e05945ce571e3e3fb543) contrib/gce: updated GCE docs with new cloud dns config and commas
 - [`db4efa7`](https://github.com/deis/deis/commit/db4efa77681ccb7d3b0e0b1706500c02c61abdf4) deisctl: update mesos code to use io.Writer
 - [`7ad2a47`](https://github.com/deis/deis/commit/7ad2a4740b626622efefd822178833ce9e37569f) deisctl: Switch from channels to io.Writers.
 - [`1755334`](https://github.com/deis/deis/commit/17553346e6dfedcad52b59f25f22e72204b4c765) client-go: defer deferred closing of body until after error check.
 - [`9bde136`](https://github.com/deis/deis/commit/9bde1363ecc6c13ccb9783361a8526a1aa228d94) Makefiles: silence static binary check
 - [`eae5d67`](https://github.com/deis/deis/commit/eae5d6758ba0e1279ae9bfd9bc91cc949b07bfea) logger: use one ticker for all time-related operations
 - [`d01053c`](https://github.com/deis/deis/commit/d01053c57e9be39b4ad06752cebe00918de7c288) commit-msg: ignore comments in line count
 - [`322fa7e`](https://github.com/deis/deis/commit/322fa7eb9bed2fad22a67be745c647024be6644b) contrib: warn if commit subject is too long
 - [`1bc6d2a`](https://github.com/deis/deis/commit/1bc6d2a0f9d8eaa9bde97285fe5c38f475bd5d83) controller: allow non superusers to manage their keys
 - [`f900d8d`](https://github.com/deis/deis/commit/f900d8d28c7bb9c990ac868a78a891edf47eded5) router: add set-misc module to nginx build
 - [`4016cc6`](https://github.com/deis/deis/commit/4016cc62899d578220effb056b03b90bc20f5b8d) (all): ensure container removal with --rm

#### Documentation

 - [`9752e08`](https://github.com/deis/deis/commit/9752e089855d4af56d061333eb09f55bdf08b40d) managing_deis: add section for registrationMode
 - [`d8aa0ca`](https://github.com/deis/deis/commit/d8aa0ca448c3075ebf2f08cce4f2f328107d3612) platform: document how to isolate etcd
 - [`1f9914c`](https://github.com/deis/deis/commit/1f9914cacc97be3d48aa906dd9b4891905e94766) reference: fix mistakes in perms section of api docs
 - [`1c76e6f`](https://github.com/deis/deis/commit/1c76e6f5358603b69824af5967c34202551f8995) MAINTAINERS: add Kent Rancourt as a core maintainer
 - [`78c1663`](https://github.com/deis/deis/commit/78c166335461eaf7ff6c22028c726c2443c0a2ef) swarm: add a basic README.md
 - [`51e2e63`](https://github.com/deis/deis/commit/51e2e63cbb6bb5b899fc8aa50bccf9e6153e8f58) concepts: remove unlinked 'the'
 - [`2efd486`](https://github.com/deis/deis/commit/2efd4864e627ec5324d01c1621d6d6d40e0c24aa) requirements: remove old warning about clang
 - [`95932cd`](https://github.com/deis/deis/commit/95932cda678bd81a433d78458c88f3945eb0dd4e) contributing: add section on Design Documents
 - [`23690c3`](https://github.com/deis/deis/commit/23690c32441ba10b126a017008f2fe165cdb131a) contributing: add commit-message hook
 - [`74e92d3`](https://github.com/deis/deis/commit/74e92d38c94a3026d38caf4f3b47e4190704eafa) reference: Add procfile as a option parameter to builds
 - [`5ebdc99`](https://github.com/deis/deis/commit/5ebdc995e65dd8b03106a21e6ff1a1e35527662d) roadmap: fix link to Kubernetes preview

#### Maintenance

 - [`be8c7b9`](https://github.com/deis/deis/commit/be8c7b97221bf878e6e17216d7ac4cefa3bf29fa) tests: update vendored docker code to 1.5.0
 - [`8d238ba`](https://github.com/deis/deis/commit/8d238badbf3b6cd96eb4eedfbc95e7ea60b788ee) deisctl: update vendored fleet code to 0.9.2
 - [`d332f73`](https://github.com/deis/deis/commit/d332f73ea523b4f0b276b3dd39e77fd979d96aaa) (all): use etcd2 inside a container, extended dance remix
 - [`f2eecf4`](https://github.com/deis/deis/commit/f2eecf4d8305302ec855911ec6cbc35f3c034eb7) contrib/ec2: rename ec2 to aws
 - [`420ee7c`](https://github.com/deis/deis/commit/420ee7c47c04b91a08a277956e233939f7569969) (all): update lists of supported providers
 - [`2a7c768`](https://github.com/deis/deis/commit/2a7c7684efdbd940b23dc2811970840b3910b346) contrib: drop Rackspace support

### v1.7.3 -> v1.8.0

#### Features

 - [`59ff8d2`](https://github.com/deis/deis/commit/59ff8d2e979e4eae7feee06c63da838e4587f259) publisher: check app config for HTTP health check
 - [`247276e`](https://github.com/deis/deis/commit/247276e7ef9bb1b9e4b00b13a75a88c3db387991) deisctl: support stateless control plane without Ceph
 - [`ea003a2`](https://github.com/deis/deis/commit/ea003a2a0e8af439de9d13730f2e23d0522deb80) contrib/linode: add Linode provision scripts and docs
 - [`11a6698`](https://github.com/deis/deis/commit/11a6698af073242bdde24e3850a43e867e5894e8) router: enable HSTS when enforceHTTPS is set
 - [`a4f0bbc`](https://github.com/deis/deis/commit/a4f0bbc4a27beec594a9f0cd6c9da7fae0bdedc4) router: add Diffie-Hellman parameter for DHE ciphersuites
 - [`fc8d2a6`](https://github.com/deis/deis/commit/fc8d2a6a51f652ffb59929598ae857360ebbdbfb) router: read list of preferred ciphers
 - [`02364d8`](https://github.com/deis/deis/commit/02364d879802dd7cc3650cc1e34abad85c765b34) router: use new reuseport option
 - [`c074db5`](https://github.com/deis/deis/commit/c074db5a034f010a97b9864994ae5f39766b0dfd) tests: retry push to dockerhub
 - [`70cd4c6`](https://github.com/deis/deis/commit/70cd4c697ddf9ab51d503f259a874eac5af4e84e) controller: Add endpoint to regenerate tokens
 - [`95d3a2e`](https://github.com/deis/deis/commit/95d3a2e8f901f36f77cd559de66359f3bf5567de) deisctl: Add deisctl ssh

#### Fixes

 - [`b6e7c2b`](https://github.com/deis/deis/commit/b6e7c2b92fa49314fd984dd9f73657681ab5616d) controller: revert HTTP healthcheck feature
 - [`4ca4f73`](https://github.com/deis/deis/commit/4ca4f738b043dea40979b8d90083abd79b5eba1a) publisher: display http error if connection failed
 - [`e9d4ef1`](https://github.com/deis/deis/commit/e9d4ef1e26dc9d5b340ed3118b8e32c0e33a53d8) builder: ignore etcdctl stdout in cleanup script
 - [`6943459`](https://github.com/deis/deis/commit/6943459a4e2d1ed7dd4163575bb12f8c50147fd6) builder: set private key perms to 0600
 - [`5faa404`](https://github.com/deis/deis/commit/5faa404250757e9eb6dd27cf36bcfae60eafe92c) controller: test that RC is available before accesing it
 - [`3250dee`](https://github.com/deis/deis/commit/3250dee80bde2e896d113d161942149f7112aeb0) tests: check https://$DEV_REGISTRY
 - [`4e2204e`](https://github.com/deis/deis/commit/4e2204ee303f569dda81cfc3d1e0f5e2a9d964cd) logger: reduce severity of log failure
 - [`b7cf655`](https://github.com/deis/deis/commit/b7cf6555356d6d3f6a080985eb2151c26f269e33) builder: only build slugbuilder/slugrunner when necessary
 - [`07db8f0`](https://github.com/deis/deis/commit/07db8f08be2257ab97c860cd6326e5273713be24) registry: run pkill as registry user, removing "-u"
 - [`47606ac`](https://github.com/deis/deis/commit/47606acd09256387abf4588880f40bc3d5079a92) docs: mock out python ldap modules
 - [`ba9899e`](https://github.com/deis/deis/commit/ba9899e714b0e60073aa72161c2dd319b66d4c62) docs: Correct and clarify description of builder component
 - [`428ecc5`](https://github.com/deis/deis/commit/428ecc5990acff7741372026c74a3e8bdd93c35f) client+deisctl: remove colors from installer messages
 - [`3f79281`](https://github.com/deis/deis/commit/3f79281212bed10433a70cdb67c647bc3c1bc94d) builder: generate key on boot
 - [`5f7e4fd`](https://github.com/deis/deis/commit/5f7e4fd5678fc4aef4d90d9340004d5a075444fc) controller: Roll back scale on failure
 - [`94a7562`](https://github.com/deis/deis/commit/94a7562abed99206fe59d88adcc56f59c0745872) builder: add "z" flag for compressed tar archives
 - [`59237fe`](https://github.com/deis/deis/commit/59237fe343cdb10798c4677fad6aace32581e9c1) deisctl: return 0 when only printing version
 - [`70cff5b`](https://github.com/deis/deis/commit/70cff5bf2375ceaf27e4c13b5e7c6fa8e3d16ff3) deisctl: avoid blocking if output is empty
 - [`5765871`](https://github.com/deis/deis/commit/5765871615ded4a43c4ab35689033df1a4fac9b8) tests: fix assertions on default beverage
- [`c6c65d1`](https://github.com/deis/deis/commit/c6c65d1461247bb20516355a3daae8e6a13d31be) router: set_random belongs in a `location` block
 - [`daa4e69`](https://github.com/deis/deis/commit/daa4e69cd77add58242c130919d724bd130a0dca) router: fix markdown table rendering
 - [`1cc6021`](https://github.com/deis/deis/commit/1cc6021bbe323171f21bbc447b0dfa717a34bdcf) client: Bump pyinstaller to fix pyinstaller/pyinstaller#152
 - [`bbe2a29`](https://github.com/deis/deis/commit/bbe2a29d168787367eec48ddd98f1bbdf9eb4b45) database: avoid pbr runtime error by pinning oslo-config

#### Documentation

 - [`26f6a3c`](https://github.com/deis/deis/commit/26f6a3cc5db5b619fde47b162d2cda27f4e61e03) managing_deis: rename stateless platform docs
 - [`8e44f76`](https://github.com/deis/deis/commit/8e44f76eb73abeb66873b2b77fdc1e4b21095624) roadmap: add July 2015 planning meeting
 - [`d6bea8d`](https://github.com/deis/deis/commit/d6bea8da8e8fdf79eba5ec5a5fd873394a3e8c67) reference: Fix response of /v1/apps/<app id>/run endpoint
 - [`4931308`](https://github.com/deis/deis/commit/4931308d0c10fad49faacd22ba6dc329b6d81b7c) reference: Fix GET /v1/keys/ response
 - [`bedd3ef`](https://github.com/deis/deis/commit/bedd3efee4c232660ccfde3dc80d5e4b59aec16d) reference: fix config list and set examples
 - [`8affe35`](https://github.com/deis/deis/commit/8affe35b2f9739bb3ac047729f635b950f8ebd8a) (all): update references to the new base image, alpine
 - [`3449a9d`](https://github.com/deis/deis/commit/3449a9d6d2dd43bb0e7120237d7eaec54cf28942) reference: add missing packages to autodoc
 - [`fa17b49`](https://github.com/deis/deis/commit/fa17b494315ea80aa51be77f195e55ed944271bb) router+upgrading: clarify EC2 PROXY protocol use
 - [`dcbb706`](https://github.com/deis/deis/commit/dcbb70656a77809179d96420133c07effdb919ae) installing_deis: include Deis Pro in AWS docs
 - [`714a669`](https://github.com/deis/deis/commit/714a6691eb142c12e4a93561d4073a77ea883e1b) installing_deis: update DO supported regions
 - [`3cc3baf`](https://github.com/deis/deis/commit/3cc3baf0fe7084bfe11b1dad47dc377879a14a2d) backing_up_data: Add community tools for reference
 - [`00e47e6`](https://github.com/deis/deis/commit/00e47e65149c774e6521d16825570ce981b17368) manage_deis: - add SPM Performance Monitoring
 - [`a135bbc`](https://github.com/deis/deis/commit/a135bbcba8ec7c0d15f34d474128b48a5080a9e2) installing_deis: add Engine Yard
 - [`3d96f64`](https://github.com/deis/deis/commit/3d96f64c4025f77eb9d3b0774eb3b00113c0a771) contrib: add link to Deis Backup and Restore
 - [`ddee95f`](https://github.com/deis/deis/commit/ddee95f13f8653a59a52c9a67851390ef87c3123) installing_deis: etcd doesn't require an odd # of nodes
 - [`1536bfc`](https://github.com/deis/deis/commit/1536bfcccb1853d56b151b6e58fa9568ecdbb72f) roadmap: update release checklist
 - [`467a07d`](https://github.com/deis/deis/commit/467a07d31464336c39c8ae4729c4af23e913bbbe) hacking, aws: Update to latest installation process.

#### Maintenance

 - [`075cebe`](https://github.com/deis/deis/commit/075cebe468778bb2c66a461dc1065ceb3c095ced) (all): revert CoreOS to 647.2.0
 - [`b02d089`](https://github.com/deis/deis/commit/b02d08931c94334a3d41f126ed7322c7fb50940d) registry: use updated registry V1 fork
 - [`990b123`](https://github.com/deis/deis/commit/990b12365d94598c347ec0ed4d01ef1c2b5a7dea) database: update wal-e to 0.8.1
 - [`d09f128`](https://github.com/deis/deis/commit/d09f128f4a9634dd4ec85df77004d51e53afc6fb) (all): update confd to v0.10.0
 - [`322658f`](https://github.com/deis/deis/commit/322658f07082d16b984703faa2f92085734001e7) contrib/ec2: add m4 instance types
 - [`8b60afb`](https://github.com/deis/deis/commit/8b60afb5f87a46e648dcd1029eb6689feb468bd4) (all): bump CoreOS to 681.2.0
 - [`ab4544b`](https://github.com/deis/deis/commit/ab4544b4918ff4b87f864086f8fe8ae0dee22e89) client+controller: update flake8 to 2.4.1
 - [`7dfa369`](https://github.com/deis/deis/commit/7dfa3694ff715a18441be0eba2c318ca9f29d3a3) controller: update psycopg2 to 2.6.1
 - [`909df2d`](https://github.com/deis/deis/commit/909df2d495aab768eb9e17e5aa1019abf7804ce7) (all): update python installer tool pip to 7.0.3
 - [`27ea62f`](https://github.com/deis/deis/commit/27ea62f12d1a1f1166bb4b1acd80008130ef80b6) (all): bump CoreOS to 681.0.0

### v1.7.2 -> v1.7.3

#### Fixes

 - [`87d2ba3`](https://github.com/deis/deis/commit/87d2ba3caaafe89b911589ff149a932f31e31cf4) builder: call ls rather than using a wildcard

### v1.7.1 -> v1.7.2

#### Fixes

- [`c73c8f7`](https://github.com/deis/deis/commit/c73c8f7e0fc96576f90a6cd4c2a66d08ba2dbae9) logger: adjust publish and ttl times to seconds

### v1.7.0 -> v1.7.1

#### Fixes

 - [`e8df401`](https://github.com/deis/deis/commit/e8df4010f52378b2ebd6195581080d82f2b3acf9) deisctl: Add extension when refreshing units

### v1.6.1 -> v1.7.0

#### Features

 - [`666c7fc`](https://github.com/deis/deis/commit/666c7fc53b7b986eddfe8e8ff134c1c72b508585) router: allow customizing controller subdomain
 - [`d51816e`](https://github.com/deis/deis/commit/d51816ed7b01e70c27d2734a39ff28150ed63c79) client: add --username flag to auth:passwd
 - [`917973d`](https://github.com/deis/deis/commit/917973d3973d8651fabcc5441e68bf047a18263e) deisctl: add config rm
 - [`0f87b9d`](https://github.com/deis/deis/commit/0f87b9dd1d04d7fd9726849ff6215a433de22c95) contrib: optional graceful shutdown
 - [`a8e4fe5`](https://github.com/deis/deis/commit/a8e4fe50e3fe0960cddb7d3e7c916c2ce223f8de) contrib/azure/azure-coreos-cluster: automatically create discovery url
 - [`f189dde`](https://github.com/deis/deis/commit/f189dde5097691d2f0c5a82cc93b7a67ff3c1114) controller: more info with bad scale operation
 - [`ec4d494`](https://github.com/deis/deis/commit/ec4d4946ddd844a4e4ee6a54a9e6ef8c196cf9c4) tests: keep test instances with SKIP_CLEANUP=true
 - [`32d848f`](https://github.com/deis/deis/commit/32d848f2f757b2dbee921b7c1e9dcc7d8395266b) controller: Deprecate X prefixed headers
 - [`33bc75f`](https://github.com/deis/deis/commit/33bc75f09f2b2570d31386a8b093fadb516b23a2) controller: expand random namespace for apps
 - [`8207364`](https://github.com/deis/deis/commit/820736474e1466d34674491b0f51ec8c236e4174) publisher: add pprof to enable profiling
 - [`ca22cf2`](https://github.com/deis/deis/commit/ca22cf2126d4db1928b7f88cf0219dfd8ba55b73) contrib/ec2: add support for custom volume sizes
 - [`48c965d`](https://github.com/deis/deis/commit/48c965d28ab3f53b96f4d191c3a237e1fecf8c0f) contrib: AWS cli profile param
 - [`eb85b8a`](https://github.com/deis/deis/commit/eb85b8a36c5cec74e4836b0ba5a783fce0ae8770) publisher: expose config as flags

#### Fixes

 - [`3290a58`](https://github.com/deis/deis/commit/3290a58df935dd0eb3732a663248f543243fa075) client: Fix attribute errors on windows
 - [`c2cffb9`](https://github.com/deis/deis/commit/c2cffb969a1edaaa456033a05b7a6958eea02bad) logger: poll etcd rather than use etcd watch
 - [`d8f012d`](https://github.com/deis/deis/commit/d8f012d844e83ed9d9c24b36a0c56dd51a8b9d0e) (all): check for key already exists errors only
 - [`00ed54c`](https://github.com/deis/deis/commit/00ed54c33c3cf41fde80ff3b22e878b68c475c4b) router: avoid consistently hashing on empty string
 - [`90fae7d`](https://github.com/deis/deis/commit/90fae7de0a472739ac08a9c27a0b4968a4e19e1a) tests/fixtures: use env var for test-etcd version
 - [`1ca5016`](https://github.com/deis/deis/commit/1ca50167825b68be66131e9773da48269f0cf72d) logger: gracefully handle a nil return from etcd watch
 - [`34cce83`](https://github.com/deis/deis/commit/34cce8351f928ad01522ecb4548cdd9e1c35eb33) logspout: use $IMAGE_PREFIX like other components
 - [`d486c84`](https://github.com/deis/deis/commit/d486c84c8ebf1c23ae8fbf8d98a1a73e44f97e85) database: work around wal-e error by pinning pbr lib
 - [`cc68a62`](https://github.com/deis/deis/commit/cc68a621563c6f30fc69ff996505226e5659466d) client: require the tag for builds:create
 - [`cb7fe2d`](https://github.com/deis/deis/commit/cb7fe2da792f9742b2708275a55efd4200f47fcc) (all): docker.service depends on its mount
 - [`b82c5fe`](https://github.com/deis/deis/commit/b82c5fe7968caf6cea223c69a055605ddf0b7b72) deisctl: hide random "deis-" services from "deisctl list"
 - [`d5a1233`](https://github.com/deis/deis/commit/d5a123349b56a70d91b92ec142d4affee88f352a) controller: merge structure with new structure
 - [`b06b85a`](https://github.com/deis/deis/commit/b06b85a474262c643c61377bd756e768420c489a) client: let requests read env vars
 - [`ebc3df4`](https://github.com/deis/deis/commit/ebc3df479a672acb4a81525ba1edd0fb4d70a6a9) tests: test custom buildpacks for example Procfile apps
 - [`bdea672`](https://github.com/deis/deis/commit/bdea672a538a33cbebbcca8b3eadb4b7a9034b7d) logger: avoid hitting etcd for each log line
 - [`7830ade`](https://github.com/deis/deis/commit/7830adea85bde3911dbccc9305f407337d1b36ce) database: run backups in the background
 - [`ba950c2`](https://github.com/deis/deis/commit/ba950c2e0ffb1cd87275d0aaaf846b4f85b36c29) (all): suppress "Key already exists" errors
 - [`21a0476`](https://github.com/deis/deis/commit/21a0476dcc03d57b57396852a3691a28411da3b7) deisctl: include swarm in refresh-units command
 - [`7e63160`](https://github.com/deis/deis/commit/7e631609c66a0f8a2bde496d20752984ac5d1be9) deis: validate the syntax of the environment file
 - [`b21d6f4`](https://github.com/deis/deis/commit/b21d6f490c6241dc523462dfdf33fad9ca2ccf00) database: remove postgres pid error from log
 - [`a222a80`](https://github.com/deis/deis/commit/a222a80970cdda33a618fb4e270d79a1f5074b54) controller: force uniqueness for default app name.
 - [`47103c6`](https://github.com/deis/deis/commit/47103c602377e4c945de15ca7ca58bc3df5155f1) builder: remove warning when docker uses the aufs driver
 - [`d89ea98`](https://github.com/deis/deis/commit/d89ea98120bab9fa5e825fc6048541abb3c67b9f) contrib/ec2: use EBS volume for etcd data
 - [`3d705a0`](https://github.com/deis/deis/commit/3d705a0f868a6adb46e0a726a1acf3e252b9fad9) nonstring-default: numeric default for webEnabled
 - [`c9d48e7`](https://github.com/deis/deis/commit/c9d48e78439783e5ebe88d308117c115b5565ee3) deisctl: prefix unit name with "deis-" before expanding
 - [`df00fb4`](https://github.com/deis/deis/commit/df00fb4e41ed1f825e6c210c637d2e4772326f64) docs: use proper deisctl syntax
 - [`368bc99`](https://github.com/deis/deis/commit/368bc99084bd0d9fe31f295c9ad5d57b49bd5e0c) contrib/azure/azure-coreos-cluster: extend load balancer timeout
 - [`6f0660c`](https://github.com/deis/deis/commit/6f0660c39ea0cac7dd8d87034372bbfc944f4d66) router: return 404 for non-existing app
 - [`1dadb72`](https://github.com/deis/deis/commit/1dadb7206ea85f532a32313473a09713fe5568a8) builder: silence pipefail errors on boot

#### Documentation

 - [`e617195`](https://github.com/deis/deis/commit/e617195295e5c6a786bb8ccf14ed73007188f855) roadmap: update roadmap
 - [`445e50d`](https://github.com/deis/deis/commit/445e50dd3613af24ff6cb7797b00a56d8b797469) MAINTAINERS: add Matt Butcher as a core maintainer
 - [`9688565`](https://github.com/deis/deis/commit/96885656f873700155f912af199e563109990eb5) managing_deis: Add notes on persisting CoreOS upgrades
 - [`ec498c6`](https://github.com/deis/deis/commit/ec498c6c7ea0f7d34473878388286c32c03d90d5) managing_deis: update backup/restore docs
 - [`d081858`](https://github.com/deis/deis/commit/d081858b5e3be0151061e0d0f2895a615b98bf37) configure-dns: drop references to local.deisapp.com
 - [`ab2df4a`](https://github.com/deis/deis/commit/ab2df4a1d32426baf1cb0e23687b4d0239d64188) roadmap: link to maintainers.md
 - [`7dd6d38`](https://github.com/deis/deis/commit/7dd6d38b4e431a8726aa84b59fd141577ff506a7) config-application: add a config:push / pull example
 - [`4f54269`](https://github.com/deis/deis/commit/4f542694b2f69c8128fd8b578c268c2f7a68a9a3) customizing_deis: clarify usage of registrationMode
 - [`86cccf4`](https://github.com/deis/deis/commit/86cccf41b6ff7acabb43b24b86b1d1d27cad0ffd) components: clarify that deis-cache is optional

#### Maintenance

 - [`56c2bab`](https://github.com/deis/deis/commit/56c2bab1d489435023fe2f9e7a05f3a02a09243a) (all): update etcdctl to v0.4.9
 - [`801ae74`](https://github.com/deis/deis/commit/801ae74455232c1cd8544c328949d777e885d628) (all): bump CoreOS to 647.2.0
 - [`f824206`](https://github.com/deis/deis/commit/f824206c78fec01c94b43b6d14cc854288d44a03) (all): bump CoreOS to 647.0.0
 - [`c0e662e`](https://github.com/deis/deis/commit/c0e662e0247f56b2fd420fdec4e4476d3f26c9e7) controller: remove $PORT logic from fleet

### v1.6.0 -> v1.6.1

#### Fixes

 - [`13bd5f2`](https://github.com/deis/deis/commit/13bd5f2c70054b7ea9c45d802b44622d578d21b3) client: pin cryptography library at 0.8.2
 - [`5e09bff`](https://github.com/deis/deis/commit/5e09bff5f7fd061af0ca6a07b2417d213d322517) store: fix shared etcd key defaults
 - [`8349d06`](https://github.com/deis/deis/commit/8349d06fe8eeee5a0c29fdeced263bc515994b87) router: check if there is certificates to generate
 - [`5dd03e8`](https://github.com/deis/deis/commit/5dd03e8705d82877ec731cd69df1b77f481921be) contrib/coreos: remove custom clock sync logic
 - [`7c8c60b`](https://github.com/deis/deis/commit/7c8c60ba71a344c396cb66c54fb0cedd4792065d) client: add requirements.txt to pypi distribution
 - [`575f68d`](https://github.com/deis/deis/commit/575f68dfe3c4a14be0add9d28bfcaeb7af29047d) store: lower number of placement groups

### v1.5.2 -> v1.6.0

#### Features

 - [`f08e9b0`](https://github.com/deis/deis/commit/f08e9b0a87ae8fcd2b5d4dd208e75b8269143243) controller: add swarm scheduler tech preview
 - [`b7cb0eb`](https://github.com/deis/deis/commit/b7cb0eb6f0c44629fd32165befcedd32e52a565b) contrib/digitalocean: enable ams2 and sfo1 datacenters
 - [`20a3eeb`](https://github.com/deis/deis/commit/20a3eeb0927daf3f0e1d5ba6a68ad81a9d25ff28) Makefile: Add a command to start a local postgresql container for running unit tests
 - [`d2122e2`](https://github.com/deis/deis/commit/d2122e2706529daf5c23d53be54f2fba6edeeda9) controller: add ps:restart command
 - [`958ad28`](https://github.com/deis/deis/commit/958ad282b4172e7d582e1225b6ebeb6b6edbe89a) router: nginx 1.9.0. Remove third party tcp module
 - [`3182d1c`](https://github.com/deis/deis/commit/3182d1cd4ca7af9081f9fd56b55c4efc02482427) controller: allow user to scale gunicorn worker processes
 - [`97ce414`](https://github.com/deis/deis/commit/97ce4148e2967deff13887c2ffdef4b28c5d3382) router: customize error_log level via etcd
 - [`4f50645`](https://github.com/deis/deis/commit/4f506451a747a59356649ddd8a08e5198dd99cb7) controller: support SAN certificates
 - [`aaa2f8e`](https://github.com/deis/deis/commit/aaa2f8eae2fbee577050b18e42131cb3b3d63a50) logger: expose config as flags
 - [`603c537`](https://github.com/deis/deis/commit/603c537ce9c1b88430a06d0bcb1309ff37b92fcc) controller: add a option to limit the number of log lines
 - [`c32b775`](https://github.com/deis/deis/commit/c32b7757477caa3b68064052ab4a60d00e219fb8) client: add optional path to config:push
 - [`ebcc410`](https://github.com/deis/deis/commit/ebcc41003d760e96c42d0c47fa3d0c2d4a478d02) controller: Allow adminOnly registration
 - [`a9019be`](https://github.com/deis/deis/commit/a9019bef4397aae81ce4d044099ff4b952e91733) contrib: list Melano, an example F#/Suave app for Deis
 - [`a37d3d5`](https://github.com/deis/deis/commit/a37d3d5e4496238c84681a04e4c158777e598f62) controller: restore coverage.py report for unit tests
 - [`fd3fc8e`](https://github.com/deis/deis/commit/fd3fc8e392e0abc361b77c51d1c5c23ebe4f0f3a) controller: add users:list endpoint
 - [`4070755`](https://github.com/deis/deis/commit/4070755d502729e3a247c80176cb2d6e15908f65) contrib/azure/azure-coreos-cluster: add an https enddpoint on c… …reation
 - [`ef3b89c`](https://github.com/deis/deis/commit/ef3b89c97e7be8fbb6c41f4f2bb4607b9ea3c2de) client: check controller before attempting to register. Fixes #3224

#### Fixes

 - [`e6a078b`](https://github.com/deis/deis/commit/e6a078b4119887b4de73e8a6d04291b7cdeb823b) contrib: remove message to exporting DEISCTL_TUNNEL
 - [`b3a352d`](https://github.com/deis/deis/commit/b3a352d8722219f04daccc5c262c0d711d88c098) registry: bump etcd ttl to 20 seconds
 - [`ea76828`](https://github.com/deis/deis/commit/ea76828bd95663748aa84b53a428392c56c83f2c) test: obtain logs from deis-registry
 - [`12e79a4`](https://github.com/deis/deis/commit/12e79a4415173605d0af06bd3b7f402efd8d72d5) user-data: update fleetd 0.92 to /opt so it survives reboots
 - [`2f4b965`](https://github.com/deis/deis/commit/2f4b9656fd25936769ebfac087142ef421a00cb8) store: change template to return ip address instead of :6789
 - [`3103d0f`](https://github.com/deis/deis/commit/3103d0f9f14438000a0f8de2e5d60b20a3f2c654) store: install lsb-release package
 - [`ded90a8`](https://github.com/deis/deis/commit/ded90a8417cfdc19cad56d2a5056b911d30ce5a2) client: Fix spacing of releases list
 - [`efd02f8`](https://github.com/deis/deis/commit/efd02f842e19c42e2c56fc71f171ba18ce87a294) registry: retrieve bucket name from etcd
 - [`e8ec005`](https://github.com/deis/deis/commit/e8ec005734078234489648a87a151a8706859f34) client: use pyOpenSSL for improved security features
 - [`93b8ce8`](https://github.com/deis/deis/commit/93b8ce8007422b4ebb2033031bd61af7f233f415) router: Unify timeout values
 - [`a4bf040`](https://github.com/deis/deis/commit/a4bf040b5982dee6b92bfd4cfdfb3d57f30cdba2) router: include deis.conf if no match with an SSL cert
 - [`1df8eea`](https://github.com/deis/deis/commit/1df8eeaaf0f563d758167c454b2a8f9de612122f) controller: allow "*" wildcard in cert REST URLs
 - [`5c2ae9c`](https://github.com/deis/deis/commit/5c2ae9c4e12ab43e4ba5193aa55cc8ab8be54376) controller: use key fingerprint instead of id for uniqueness
 - [`7f6099c`](https://github.com/deis/deis/commit/7f6099c865c7a48f88b435cfad01f435251208e2) controller: return the correct domain from get_object(
 - [`86b3b09`](https://github.com/deis/deis/commit/86b3b09324e3a9fe16d565f385908038533cdb73) contrib/ec2: improve timeout handling
 - [`02b78a2`](https://github.com/deis/deis/commit/02b78a258c93d1caa149be2dbf313a097552485f) router: write out only if cert matches the path
 - [`8e960ca`](https://github.com/deis/deis/commit/8e960ca203bd5d8c48e129c4172dd2b13cdc4439) (all): use "confd --interval 5" instead of "--watch"
 - [`ce21c3b`](https://github.com/deis/deis/commit/ce21c3bfa1e71a140e9f5e259e37b1a97de5c54d) controller:  fix regex for deis limits
 - [`fa48fd2`](https://github.com/deis/deis/commit/fa48fd21e9aed44f7c50c43211b96dfa4c28f21e) controller:  remove domain
 - [`49d9002`](https://github.com/deis/deis/commit/49d9002e1cd2b1b0cac60c1aa6d3926c994fb3fb) builder: exit 1 when gitreceive captures a build
 - [`bd8b2df`](https://github.com/deis/deis/commit/bd8b2dff4b03efa66d77f22078f6d2ccceceb00f) router: use nginx $host in HTTPS redirects
 - [`def0d96`](https://github.com/deis/deis/commit/def0d9607e58cada0555a99ce86dd1e682954c8d) 'contrib/azure/azure-coreos-cluster': make this not found error more user friendly
 - [`e153834`](https://github.com/deis/deis/commit/e1538349fad1be04d906ecb04ce2f1674b3abfed) database+registry: check for S3 bucket name before creating it

#### Documentation

 - [`3fa5df5`](https://github.com/deis/deis/commit/3fa5df5de89262f9321cec63e1411193fb0f6787) contrib: add link to Docker S3 Cleaner
 - [`ee9f00e`](https://github.com/deis/deis/commit/ee9f00ed1a78f9b7ecf90f919058f5d1d0d6447d) using_deis: explain difference between web and cmd
 - [`54757a0`](https://github.com/deis/deis/commit/54757a08e5f714789c0ead8307736849e2412136) contrib: organize community contributions
 - [`6fe16b1`](https://github.com/deis/deis/commit/6fe16b1f90b6902cd842c8d3567af60d135c7d50) installing_deis/aws: update AWS provisioning example output
 - [`b48ced3`](https://github.com/deis/deis/commit/b48ced3e3ad2c964634fcefd5f8cb56f712677a7) troubleshooting_deis: use private key on ssh -i
 - [`ce4f8c3`](https://github.com/deis/deis/commit/ce4f8c3b97e7f6c6289d80f4010096776c78bf8a) managing_deis: add upgrade Deis clients for in-place upgrade

#### Maintenance

 - [`f4c2a24`](https://github.com/deis/deis/commit/f4c2a241530522d0d612fe3448bf9309d2c7927f) builder: migrate to heroku's cedar stack
 - [`84cd7d0`](https://github.com/deis/deis/commit/84cd7d0d6e88725ffdbeb560f8c91d2c28f9d25a) database: update wal-e to v0.8.0 + busybox fix
 - [`4b6cc67`](https://github.com/deis/deis/commit/4b6cc67dfa6874c4cb5206bbda6b38b3f4bad0d8) contrib/rackspace: bump CoreOS to 633.1.0
 - [`053eb4e`](https://github.com/deis/deis/commit/053eb4e7a8b94b5f5abd445193aab61ddeef3662) router: update nginx to 1.8
 - [`858d244`](https://github.com/deis/deis/commit/858d244be6df94286c58fcbeb8413197705b4573) contrib/coreos: remove etcd LimitNOFILE
 - [`0a8c9e3`](https://github.com/deis/deis/commit/0a8c9e398a0695fc9447ddd8699ad723a7df19bb) (all): bump CoreOS to 633.1.0
 - [`6702716`](https://github.com/deis/deis/commit/670271697c236743dfcfd3e8d4f85f1fa6922d35) (all): bump confd to v0.9.0
 - [`82cf84c`](https://github.com/deis/deis/commit/82cf84c6c8c7bd6c898d98df41d8b325620106e2) controller: update docker-py to 1.1.0
 - [`1cdf41e`](https://github.com/deis/deis/commit/1cdf41e929d12f8f98c5ac638b1294c0243a9a9b) (all): update pip installer tool to 6.1.1
 - [`d120816`](https://github.com/deis/deis/commit/d1208162f232b1f7309b1dc4f2b98bcb76376233) contrib/ec2: introducing c4 instance types
 - [`da4408b`](https://github.com/deis/deis/commit/da4408b87eafce9c5d00e63bf0b716add1de3080) store: bump Ceph to "hammer" release
 - [`fb081f8`](https://github.com/deis/deis/commit/fb081f86cd1d9b5671bc489d71a64123d8d0fddd) client: update python-dateutil to 2.4.2
 - [`ae7adc1`](https://github.com/deis/deis/commit/ae7adc16f21a0eff5685d86627fdb8fe6c577d79) logspout: switch to busybox

### v1.5.1 -> v1.5.2

#### Fixes

 - [`a4bf040`](https://github.com/deis/deis/commit/a4bf040b5982dee6b92bfd4cfdfb3d57f30cdba2) router: include deis.conf if no match with an SSL cert
 - [`1df8eea`](https://github.com/deis/deis/commit/1df8eeaaf0f563d758167c454b2a8f9de612122f) controller: allow "*" wildcard in cert REST URLs
 - [`7f6099c`](https://github.com/deis/deis/commit/7f6099c865c7a48f88b435cfad01f435251208e2) controller: return the correct domain from get_object
 - [`02b78a2`](https://github.com/deis/deis/commit/02b78a258c93d1caa149be2dbf313a097552485f) router: write out only if cert matches the path
 - [`8e960ca`](https://github.com/deis/deis/commit/8e960ca203bd5d8c48e129c4172dd2b13cdc4439) (all): use "confd --interval 5" instead of "--watch"

### v1.5.0 -> v1.5.1

#### Fixes

 - [`0866951`](https://github.com/deis/deis/commit/0866951334e3affaa9d04ae8f89f95bfe2590682) builder: exit 1 when gitreceive captures a build
 - [`64c0990`](https://github.com/deis/deis/commit/64c0990e8cf1d71a45e2361f9d389abdfc53f0f0) router: use nginx $host in HTTPS redirects

#### Maintenance

 - [`f66ae3b`](https://github.com/deis/deis/commit/f66ae3badfd9f330a7b0caa7fd82aaef9521b9e8) release: update version to v1.5.1

### v1.4.1 -> v1.5.0

#### Features

 - [`7c11861`](https://github.com/deis/deis/commit/7c11861dde635ff53136411b0a7c8c58156ed1ab) tests: add cert integration tests
 - [`f376230`](https://github.com/deis/deis/commit/f3762305947e176894ffa4da556a48662607575f) tests: add domain integration tests
 - [`9114020`](https://github.com/deis/deis/commit/91140200b4b27141b889a9f9465f6e7663654a0c) controller: add deis certs
 - [`27379c7`](https://github.com/deis/deis/commit/27379c71e3f26392a9559e939010bb2965cc7f73) client: add client extensions
 - [`b5644c3`](https://github.com/deis/deis/commit/b5644c3e7261407c3904d416b40d9da0cd388c8c) tests: use m3.medium instances for ec2 integration tests
 - [`5a5ab0a`](https://github.com/deis/deis/commit/5a5ab0a4ef95cf2d25b79ea9eee91f11e5ac7df4) (all): use confd --watch
 - [`1fd7b59`](https://github.com/deis/deis/commit/1fd7b592aaf7f3a4b0eba1ad905a85e39b7335d8) client: add deis config:push
 - [`2f29bcf`](https://github.com/deis/deis/commit/2f29bcfa7609111c482fd0ce7cbb07dbc1e74f09) client: Add multiple profile support to the client
 - [`e32d008`](https://github.com/deis/deis/commit/e32d0087f02bd635b9ffde2e9a6968d917871be4) controller: Adding LDAP/AD auth support
 - [`8c8c94c`](https://github.com/deis/deis/commit/8c8c94c387bcc3bc767477eabfa6c3767ca69481) controller: disable swap usage if there is a memory limit"
 - [`0182052`](https://github.com/deis/deis/commit/01820529186bc3513129c6620ba3991602d27449) test: remove hardcoded names and extract journal logs from each user app
 - [`d1a33e6`](https://github.com/deis/deis/commit/d1a33e6b3090035819930b3635395ce9e4a87de2) builder: add ssh key variable to support private repositories
 - [`b07a053`](https://github.com/deis/deis/commit/b07a0534b9a0e1c24593fb26e3a49c22578ab28a) controller: print stack trace when worker is killed
 - [`380e1bf`](https://github.com/deis/deis/commit/380e1bf85b7340f10faef0fa1f8d305b32220ef9) controller: disable swap usage if there is a memory limit
 - [`5111537`](https://github.com/deis/deis/commit/51115373cb275d08cb0c8d9f969ccd6736f88fa3) client: add flags for automating auth:cancel
 - [`a8eb3ab`](https://github.com/deis/deis/commit/a8eb3ab5ae5ff157be65a908fbbe6a1d67556058) azure: add 'affinity-group' argument to service creation script
 - [`e6eb3b8`](https://github.com/deis/deis/commit/e6eb3b84477781daab593f860728972bf8a95d6e) contrib/coreos: add preseed script for Deis components
 - [`1c7ad47`](https://github.com/deis/deis/commit/1c7ad47d5629602ec77b00e5053e59639729f2d4) builder: detect and use overlay as graph driver if supported
 - [`dd39b48`](https://github.com/deis/deis/commit/dd39b483958160e488094252116863e5d0346770) docker: load overlay module
 - [`92a5549`](https://github.com/deis/deis/commit/92a5549fe5c978d74682f5ad6896d0838a4f024f) router: add optional query string affinity support
 - [`0292ff2`](https://github.com/deis/deis/commit/0292ff2b5c4a6072d6d64502f179face6c9c428d) contrib/ec2/gen-json.py: split private and public subnets
 - [`8b5a0b3`](https://github.com/deis/deis/commit/8b5a0b313639f59b85f7d965edef6901b079e732) router: set additional X-Forwarded headers
 - [`000406b`](https://github.com/deis/deis/commit/000406bd04dd1613b2dfc10d937c3491a95a9f52) Makefile: always restart dev registry
 - [`2c44c68`](https://github.com/deis/deis/commit/2c44c68998ba4588152ca9e70a945f94ecd70d7e) Vagrantfile: add forwarded ports configuration
 - [`b61b7f1`](https://github.com/deis/deis/commit/b61b7f11bc9d31ca71fc0d11dc3f16a95ebf8676) setup-node.sh: remove user intervention setting up new nodes
 - [`614a0b4`](https://github.com/deis/deis/commit/614a0b49b0a0533b753b599dfc01349b4d105143) cloud-init: temporal fix for ntpd
 - [`b320a67`](https://github.com/deis/deis/commit/b320a67da8d3cbb02cc636d195d2bd3d5912b0e7) router: avoid regex-based taxing rewrites
 - [`dc63193`](https://github.com/deis/deis/commit/dc631934b7de28deac810745b2ec02a558294c00) docs: gcutil is deprecated; convert all calls to gcloud compute
 - [`669ab53`](https://github.com/deis/deis/commit/669ab535f99a2b9a0db0570fadabda64e0825a27) tests: add route53 utility for wildcard domains
 - [`e68d1ae`](https://github.com/deis/deis/commit/e68d1ae1c6f1a4944d2da2d45c1ea2896308f26e) tests: run integration tests against EC2 clusters
 - [`f6c1c7e`](https://github.com/deis/deis/commit/f6c1c7eb01f578fcd4498ed39b23204fbaeca5b5) deis client: add ssl-verify parameter to deis register command
 - [`7455c2f`](https://github.com/deis/deis/commit/7455c2f54808115d72d193f08c38a75b4554eda7) publisher: remove application when is stoped
 - [`22e4056`](https://github.com/deis/deis/commit/22e405672e06df69f05c95830d3eaa27814ff4be) makefile: check that generated go binaries are statically linked
 - [`61c8755`](https://github.com/deis/deis/commit/61c87552be6e8bd91c295b332683163f43878b31) client: add discrete git:remote command
 - [`311f0f0`](https://github.com/deis/deis/commit/311f0f0c4e48a3e733bb220bbe92b484c002555c) contrib/ec2: loop until instances have passed health checks
 - [`3285e21`](https://github.com/deis/deis/commit/3285e21d977a75767e846f68a0d0fa8bd2cdb9b5) publisher: Check if app port is open before publishing it. Fixes 3115
 - [`f50eb36`](https://github.com/deis/deis/commit/f50eb364e0236016045af7524be9cc95b7e6da9a) cloud-init: disable automatic core dumps
 - [`571f6f5`](https://github.com/deis/deis/commit/571f6f512df2d89026ea4b66487fd3c172a9bebd) logger: add platform wide log drain - syslog

#### Fixes

 - [`6441a4f`](https://github.com/deis/deis/commit/6441a4f03e1d6a6af7b16df5e43254d2c68834fb) registry: create /deis/cache etcdctl dir
 - [`db51c61`](https://github.com/deis/deis/commit/db51c6128e19c3d8f5e4d925cd43388d5dcbce0d) deisctl: omit cache from `deisctl uninstall platform"
 - [`0e93e47`](https://github.com/deis/deis/commit/0e93e47be6aab3f4308b7a3c118c9eb9b6c57b3e) deisctl: omit cache from `deisctl stop platform"
 - [`28ad7a3`](https://github.com/deis/deis/commit/28ad7a306efa53c967bfc77d7fd97bda16498542) contrib/ec2: open port 443
 - [`7f87e8a`](https://github.com/deis/deis/commit/7f87e8adc3a47e9d9b052e4969953d642a6a31b8) tests: disable TLS for self-signed certs
 - [`7049c9b`](https://github.com/deis/deis/commit/7049c9be91f051a0d256d5aba7ce97d63af56ac3) controller: allow xip.io domains
 - [`ba14d16`](https://github.com/deis/deis/commit/ba14d16b874b8c6f1874eba38764ff243a82e992) controller: reference username in api call
 - [`a5b1b39`](https://github.com/deis/deis/commit/a5b1b39768c702fb517bcdac1b81ea47a392a7e1) docs: update example output from provision-ec2-cluster.sh
 - [`44100c3`](https://github.com/deis/deis/commit/44100c35cccd98a5ca2265a24ef087678f3114e7) docs: mention that Deis requires 3 or more machines
 - [`617416b`](https://github.com/deis/deis/commit/617416b5de5617d67c8e1112b90c08167b118d81) database: go through recovery if launched with stale data container
 - [`8357742`](https://github.com/deis/deis/commit/835774225cd27d64425f23727075833d8cbc2b44) builder: retrieve process types from slugbuilder's Procfile
 - [`b0a87b7`](https://github.com/deis/deis/commit/b0a87b7a81a030adeeb86f13afcd33ad1d910e03) docs: bracket $HOME properly in "deisctl config" example
 - [`e9f7093`](https://github.com/deis/deis/commit/e9f709377772066f39e394614ba2ec2e980853dc) builder: skip `docker rmi` if there are no images to remove
 - [`5246753`](https://github.com/deis/deis/commit/524675316c1993102a6bc3b89977f84c1f1c938f) database: stub out reload script before starting confd
 - [`a7fd7f8`](https://github.com/deis/deis/commit/a7fd7f82b6d6b4d6660f48ccf3257295700e772d) controller: read from SSH channel before checking exit status
 - [`a96b5f0`](https://github.com/deis/deis/commit/a96b5f0aca2b24a3d899a0c33a6258974ea3bd5b) contrib: preseed all images
 - [`fe091d1`](https://github.com/deis/deis/commit/fe091d16922427bdc655493a616b69c1f9c80fb2) tests: allow concurrent client tests with DEIS_PROFILE
 - [`c0815f0`](https://github.com/deis/deis/commit/c0815f03671fe83b4f3818aa95cddd6a59ce47de) docs: Fix heroku-buildpack-php url
 - [`2f9381b`](https://github.com/deis/deis/commit/2f9381bc8711a603d4e6a82a72b6baeb4064b2d3) tests: replace hard-coded test domain in ec2 script
 - [`25ebb4d`](https://github.com/deis/deis/commit/25ebb4d9ec8d53e4e3a6b3ac18eff4467e09f8d1) publisher: use TTL in services application etcd directories
 - [`5505aac`](https://github.com/deis/deis/commit/5505aacde1adc08aec122d9429c3f8a2cef1334f) registry: silence confd errors from boot
 - [`41053f5`](https://github.com/deis/deis/commit/41053f5e2f555a3443f350adb598d732b201ab4b) builder: is not possible to use ioutil.ReadAll more than once with the same reader.
 - [`dafca34`](https://github.com/deis/deis/commit/dafca346e5d1c048ce0604ba60c32b9704cfcc5b) builder: return error message as string
 - [`7954a93`](https://github.com/deis/deis/commit/7954a930732c088d4c6d27cec01e8ea4e16a4c84) controller: ignore KeyError on delete
 - [`a3fce8b`](https://github.com/deis/deis/commit/a3fce8bc82e0a7f377f2e43ab7b6bc2d82705395) client: read response from DELETE /v1/auth/cancel
 - [`ea765c8`](https://github.com/deis/deis/commit/ea765c82524d33ab9e6076f6f465c5a558e73f3e) tests: lower cyclomatic complexity of Checklist
 - [`443684f`](https://github.com/deis/deis/commit/443684ff0ad4e8fdb7c835b370d5ccacca2f4b2f) docs: missing controller requirement
 - [`415739d`](https://github.com/deis/deis/commit/415739d05e34a9dfed91d905f552582d20dc405d) contrib/azure: fix virtual network name
 - [`c61566e`](https://github.com/deis/deis/commit/c61566e8eab0b979bf0c1eefdcc824f94ab3601b) deisctl/units: eat error output for history and inspect
 - [`d50cea7`](https://github.com/deis/deis/commit/d50cea7f89a3f7ac8ebf127d914db5fe0c6fb8e1) controller: remove dead web link to logout page
 - [`0191bbd`](https://github.com/deis/deis/commit/0191bbd82c392324a223a58febcf9af5175fcf83) builder: use docker import instead with cedarish
 - [`64ed491`](https://github.com/deis/deis/commit/64ed4917c501a1bc84d7f6bb945622e746ea1b7e) tests: restore pipefail error checking to ec2 script
 - [`18ee46a`](https://github.com/deis/deis/commit/18ee46a4ac48d5907ec2e5aac30b8b7483fd0875) controller: allow multiple unit tests with random db name
 - [`cf8a825`](https://github.com/deis/deis/commit/cf8a825f62a0ac048cf0df16032f5b8edb5d83d4) user-data: use LimitNOFILE for etcd service
 - [`f934dbe`](https://github.com/deis/deis/commit/f934dbe04f27f9da46fa601c179f074c020a0df7) controller: retrieve requested domain from .get_object
 - [`1cbf311`](https://github.com/deis/deis/commit/1cbf3112fc09727ff5a0fce5e0263e1f6dd021eb) database: add requirements for wal-e build
 - [`b193700`](https://github.com/deis/deis/commit/b1937008f02785618ebad81047ab8752f772258c) controller: Update CORS whitelisted request headers.
 - [`7681a4d`](https://github.com/deis/deis/commit/7681a4d08f61cae12be6f8d3146036bd54190dae) tests: create SSH key needed for keys_test.go
 - [`5caf386`](https://github.com/deis/deis/commit/5caf3862bff3b32a34f2816b4fbd80bd2bcc2353) tests: remove vagrant-specific regression test
 - [`7a29b29`](https://github.com/deis/deis/commit/7a29b29f3e572325c3a3ca61c233ed0d02c07c45) route53-wildcard.py: poll until wildcard domain is accessible
 - [`a6c2e9c`](https://github.com/deis/deis/commit/a6c2e9c098397235db4c2c678d5028b689f3a76e) docker: wait for drive formatting to complete before mounting
 - [`821e166`](https://github.com/deis/deis/commit/821e166172bd654c61cf8ab134ed354afb2d537d) router: enforce HTTPS correctly when not behind an elb
 - [`54ffbf2`](https://github.com/deis/deis/commit/54ffbf2c984a4275b55f3e7d7cf03e6f98be88ec) router: allow customomize connection limits
 - [`80c4a59`](https://github.com/deis/deis/commit/80c4a5994cab9c751581496a6e258f0976b2cb1e) client: Older dates are incorrectly parsed as today and yesterday.
 - [`c698a98`](https://github.com/deis/deis/commit/c698a980e0e2f17eb43f5ebf316dfcc8770c3272) controller: silence redundant app release logs
 - [`9171c3a`](https://github.com/deis/deis/commit/9171c3a2642d0ea6d59aa968cbc94752d173536a) store: abort build in case of any error
 - [`ea67e5f`](https://github.com/deis/deis/commit/ea67e5f67f479d133c530a346ff41aa027fe9992) controller: use same API call to determine loadstate
 - [`be660e0`](https://github.com/deis/deis/commit/be660e01e6867f26cd1e222c3594f2a50c44ef30) publisher: use a mutex to protect against concurrent access errors
 - [`fc5213a`](https://github.com/deis/deis/commit/fc5213a2fe47a1cdf0011b4b4597f2f85c3e6516) docs: handle changes in Sphinx 1.3.0
 - [`50b6775`](https://github.com/deis/deis/commit/50b6775044ea5871b150f595dca097ae3b28a4d7) builder: silence clone on .download_buildpack(
 - [`e5b0b34`](https://github.com/deis/deis/commit/e5b0b34ccef73e9cfef1e581e81e6a04cabeeff3) controller: allow labels with trailing numbers
 - [`5214b9f`](https://github.com/deis/deis/commit/5214b9f743e5bff14709f1644dda05feffc76c15) registry: remove cache from install/start platform
 - [`4033b0c`](https://github.com/deis/deis/commit/4033b0c2d15b87fc3f0843b7e2e7eec574eea5f3) docs: Starting store-gateway needs @1.
 - [`5f395db`](https://github.com/deis/deis/commit/5f395db8edcbaeba1c1de418c8c18e4f37d7895d) builder: strip single quotes from BUILDPACK_URL"
 - [`69c3787`](https://github.com/deis/deis/commit/69c3787fc50ff8ae4a46cac94edec3749f1be5a5) builder: properly escape backticks"
 - [`68bdeac`](https://github.com/deis/deis/commit/68bdeacf4b561f4184b39368364a8dcda09b5ccf) builder: properly escape backticks in envvars"
 - [`1df43a6`](https://github.com/deis/deis/commit/1df43a6c6d1dc22e2554cd9fb0c2aa2f43fe546a) builder: error on pipefail

#### Documentation

 - [`4028745`](https://github.com/deis/deis/commit/4028745a217c44af5b97438fb95b6732347bd1fd) logger: change incorrect ref to db component
 - [`4f95a85`](https://github.com/deis/deis/commit/4f95a8518689d67dde1dd2c3db1c16677e37969c) faq: summarize difference between deis and deisctl clients
 - [`cbbae6f`](https://github.com/deis/deis/commit/cbbae6fd602cc71793525114338ddc095576a2f8) reference: SSL support for custom domains
 - [`fb858d9`](https://github.com/deis/deis/commit/fb858d92652923b748c5de68889213365cde2602) customizing_deis: add "CLI Plugins" article
 - [`fddf8c2`](https://github.com/deis/deis/commit/fddf8c208df1d674e7c80a6ab6684001abe3a501) installing: revise deisctl install to use sudo on CoreOS
 - [`9eb9472`](https://github.com/deis/deis/commit/9eb9472e17607cde82dc6ff9959f26a4bab8b01a) controller_settings: Documentation to use the LDAP Auth
 - [`e74e8c2`](https://github.com/deis/deis/commit/e74e8c274961e40eef7917148e6d6efde9cac731) contrib: Add deis-backup-service
 - [`4701cd2`](https://github.com/deis/deis/commit/4701cd292bb14123d2b577e8f51c6f0fe41755e1) testing: note that some containers don't rebuild automatically
 - [`56e5f8d`](https://github.com/deis/deis/commit/56e5f8d6a28697f026bfb1e07f91ab37c4ce1ebd) managing_deis: add production deployment recommendations
 - [`54ee4a3`](https://github.com/deis/deis/commit/54ee4a304ffa105f7fd2a2cc80bbc06e3e89f925) MAINTAINERS: add Manuel and Johannes
 - [`4b455ee`](https://github.com/deis/deis/commit/4b455ee428eeabaebcc4a4fcb31c4f58e26c470c) customizing_deis: remove unused registry secret key
 - [`407dbd2`](https://github.com/deis/deis/commit/407dbd2559a529483cec68504dc0452a91a5fd58) MAINTAINERS: add documentation describing project maintainers
 - [`90f5528`](https://github.com/deis/deis/commit/90f55287b711234b89b210c172457926bc2bc846) platform_monitoring: make cAdvisor unit more resilient
 - [`c4cf140`](https://github.com/deis/deis/commit/c4cf14026390f7541280cde970ca4dab80a5317c) managing_deis: remove note on upgrading
 - [`d8e63d3`](https://github.com/deis/deis/commit/d8e63d3d54d9c46d292d07a4bf39eadd07a31b26) reference: email is a mandatory parameter
 - [`ccdf226`](https://github.com/deis/deis/commit/ccdf2261570ea9d409e834c389f6e690c0f48ea7) installing_deis: Changing ssh key docs for vagrant users
 - [`d7e2007`](https://github.com/deis/deis/commit/d7e200721e97f41dd8803676d85a5fb8b477b426) managing_deis: remove scale note on backup/restore
 - [`170dc4e`](https://github.com/deis/deis/commit/170dc4e583f56109c561e76cda54054c35042f70) using_deis: clarify proxy for Vagrant

#### Maintenance

 - [`a46a325`](https://github.com/deis/deis/commit/a46a3257df60b7e97cab2c49b8133f9e64b8009a) release: update version to v1.5.0
 - [`3f3eba4`](https://github.com/deis/deis/commit/3f3eba478a6a3181652669af942b502fb3b1c558) contrib/coreos: upgrade fleet to 0.9.2
 - [`8b40b08`](https://github.com/deis/deis/commit/8b40b081cebd79e5735be1f3b02a069f0088f29e) controller: bump API version to v1.2.0
 - [`617416b`](https://github.com/deis/deis/commit/617416b5de5617d67c8e1112b90c08167b118d81) database: go through recovery if launched with stale data container
 - [`0a13242`](https://github.com/deis/deis/commit/0a13242af8c9aa14afdd3f35f1e86b3c1cc4f685) contrib: update Rackspace to CoreOS 607.0.0
 - [`8baab77`](https://github.com/deis/deis/commit/8baab77f6dc0a2e6f346dedf087cead9300b3811) builder: bump confd to v0.8
 - [`0f0ec25`](https://github.com/deis/deis/commit/0f0ec2503db3015efb7a431be58fe7ab2efcc41d) client+controller: update flake8 to 2.4.0
 - [`ed4a5e9`](https://github.com/deis/deis/commit/ed4a5e9ae24fe29fae89aaa4dbc8a3495b3157d0) store: bump confd to v0.8
 - [`c07cbd1`](https://github.com/deis/deis/commit/c07cbd18353304f9e72fbc3b74c50b9e47da12af) registry: bump confd to v0.8
 - [`983a26b`](https://github.com/deis/deis/commit/983a26b085838dc632196e2142773e5f5f441814) database: bump confd to v0.8
 - [`898a10d`](https://github.com/deis/deis/commit/898a10d9fbd3fb8726ae3a43a8da7f89c1e51d46) controller: remove check script
 - [`21b9d60`](https://github.com/deis/deis/commit/21b9d6000f342120bc86e0ac95be662ffb0dcf1e) controller: bump confd to v0.8
 - [`19f35f6`](https://github.com/deis/deis/commit/19f35f60c7c48255ca92f0c3a1286d31718c71b6) docs: update Sphinx to 1.3.1
 - [`7276c66`](https://github.com/deis/deis/commit/7276c661442ffb3735be9a56542f9b92a10f2714) controller: update Django to 1.6.11 security release
 - [`d72765d`](https://github.com/deis/deis/commit/d72765d340ac36449a9915ced5d0dcbb30f09bd2) builder: bump heroku-buildpack-ruby to v134
 - [`1a16562`](https://github.com/deis/deis/commit/1a165628b19d65284c1c9f883053ffa418e238b3) router: bump nginx to 1.7.10 mainline
 - [`7c96ef4`](https://github.com/deis/deis/commit/7c96ef434497d7d6aaccc472cb29ec466f43e310) builder: bump Docker to 1.5.0
 - [`93d9215`](https://github.com/deis/deis/commit/93d921542237d13a0e1f61fbc52f4e589820bfe5) (all): bump CoreOS to 607.0.0
 - [`16ef6b6`](https://github.com/deis/deis/commit/16ef6b695b42ebc8c24fa4778bddaa3a04faed88) tests: add etcd log output to integration tests
 - [`9a2b653`](https://github.com/deis/deis/commit/9a2b653d3f1932870f640c2e34f88e46aa78658a) database: update wal-e to v0.8c2
 - [`03a9db0`](https://github.com/deis/deis/commit/03a9db0350fe503dd034e6d2560b2fb64746418c) registry: remove unused secret key
 - [`12632e2`](https://github.com/deis/deis/commit/12632e26561dcacf703316e57ec2f52e0b370bb1) client: update python-dateutil to 2.4.1
 - [`2069fe1`](https://github.com/deis/deis/commit/2069fe1da7f43b50f8a70da626ee44bb5bebaa6b) controller: update gunicorn to 19.3.0
 - [`e372bc3`](https://github.com/deis/deis/commit/e372bc39ba6b96c9a6cb493d3c5734e5208d06f1) release: update version in master to v1.4.1
 - [`cca9843`](https://github.com/deis/deis/commit/cca98432310f320213c8738bb30f4c5f3dee2496) release: update version in master to v1.5.0-dev
 - [`9813e45`](https://github.com/deis/deis/commit/9813e4507c45485b878014ed0b57d10a5fe26e02) Makefile: bump dev-registry to use 0.9.1

### v1.4.0 -> v1.4.1

#### Fixes

 - [`cf61859`](https://github.com/deis/deis/commit/cf61859be532603b615c511ab9b2f6163e0d8ce6) builder: strip single quotes from BUILDPACK_URL"
 - [`45e7413`](https://github.com/deis/deis/commit/45e74138b6b88c14fbc844bbc8eb1979237601ab) builder: properly escape backticks"
 - [`0046494`](https://github.com/deis/deis/commit/004649486f26c4d3c43fd4f56a27a4e873b6f59d) builder: properly escape backticks in envvars"

#### Maintenance

 - [`6960646`](https://github.com/deis/deis/commit/696064661a70846ccc1780299a714aeff8c4999e) release: update version to v1.4.1

### v1.3.1 -> v1.4.0

#### Features

 - [`1c5ae0c`](https://github.com/deis/deis/commit/1c5ae0c9d9f02fa6e3374968dcbda4b53d1c8909) builder: allow to lock BUILDPACK_URL to a commit
 - [`17ccf8a`](https://github.com/deis/deis/commit/17ccf8a72d16f3e5beb0d10080ad39ae51203224) builder: Adding support to lock BUILDPACK_URL to a git revision
 - [`f35e3a3`](https://github.com/deis/deis/commit/f35e3a3ef6d32e280d28b5fa7b3e43e2bab61582) cloud-init: check if deisctl does not exists before install it
 - [`fb6f0a0`](https://github.com/deis/deis/commit/fb6f0a0cbc37086833e6b0aac70f799ea58c86e1) scheduler: graceful shutdown with SIGTERM
 - [`751cc6b`](https://github.com/deis/deis/commit/751cc6b79ffeea43f1db5c8a3e154a71545c7b3b) deisctl: start / stop all installed units
 - [`8a09204`](https://github.com/deis/deis/commit/8a09204b6b56e8599616024474bf1d72c9d9b0dd) router: Removed X-Deis-Upstream
 - [`cbb9fa1`](https://github.com/deis/deis/commit/cbb9fa146840edb82464097c6fe17c0966cb15e2) controller: worker processes named as "deis-controller"
 - [`82574fc`](https://github.com/deis/deis/commit/82574fc054dc6ef2a61e5934c84b3a34ec1af103) builder: add proxy support
 - [`b4cddba`](https://github.com/deis/deis/commit/b4cddba96a24e0b04e7d5f100e2b92facae8561a) router: add optional HTTPs redirect
 - [`dd5004c`](https://github.com/deis/deis/commit/dd5004cc7f36710efb8b34b44a6f646912332725) contrib/ec2: add support for internal ELBs
 - [`9cb77bf`](https://github.com/deis/deis/commit/9cb77bf609eec8bb21f169c1b3ab59048fe698b6) deis/config: remove "===" when listing configs with --oneline
 - [`c35f9e8`](https://github.com/deis/deis/commit/c35f9e883e50446ed71256d32f6c1fe10147ab1d) contrib/azure: clean up Azure docs and scripts
 - [`ad85aba`](https://github.com/deis/deis/commit/ad85aba1c1d38f7b5b2afe088d1b20a2e9909609) contrib/azure: add Azure provision scripts and docs
 - [`261f78d`](https://github.com/deis/deis/commit/261f78dc86f484480529dc03e34206c34225f27d) logspout: use custom datetime
 - [`b3c06af`](https://github.com/deis/deis/commit/b3c06af618ad6e06a7bda1b2788b33ab5dba8d8b) router: nginx status; log http_host, upstream and request times
 - [`e35206b`](https://github.com/deis/deis/commit/e35206b2c74fcc7bb172141d9cc7effefa7d267a) builder: Try shallow cloning buildpack repos
 - [`6aaa48f`](https://github.com/deis/deis/commit/6aaa48fadc649d4ab7397e450c3f2e455bb9a47d) ec2/cloudformation: Select SSD EBS volumes by default
 - [`db946a0`](https://github.com/deis/deis/commit/db946a0901f2599c669081a017c85f55fd90cd1f) store: Scalable store-gateway

#### Fixes

 - [`883c4f9`](https://github.com/deis/deis/commit/883c4f91b1897bc5acd7d4fca84752915d447fea) builder: install docker from get.docker.com
 - [`bdc4313`](https://github.com/deis/deis/commit/bdc431327884e5f082a89610e96f990652db8e40) tests: check for "git "
 - [`70ae060`](https://github.com/deis/deis/commit/70ae0606d6f7ce035342b8e2c584d994b7b26dbf) tests: add --app
 - [`0c6353f`](https://github.com/deis/deis/commit/0c6353fda1d4f22c1e6778b5d5ac65c61a543f80) builder: silence errors from initial clone
 - [`1598d72`](https://github.com/deis/deis/commit/1598d72d23f423cd877e696fe66230aade3d689b) builder: strip single quotes from BUILDPACK_URL
 - [`30b3306`](https://github.com/deis/deis/commit/30b330605cb7c490bf1cc898185f3e8370fde471) tests: remove mock-store containers after functional tests
 - [`3c462b7`](https://github.com/deis/deis/commit/3c462b7b9a67c49ed2cc9b7f8b2ab89777bb4150) controller: fixup fleet reporting failed state
 - [`9ee7525`](https://github.com/deis/deis/commit/9ee752507deba2c2181334bcfc11cd6f6c212138) controller: remove print calls
 - [`cb82f62`](https://github.com/deis/deis/commit/cb82f627d65f9a372439ce2649c1d9fbaff56d66) controller: destroy app containers in parallel
 - [`f811b7d`](https://github.com/deis/deis/commit/f811b7ddddb6af19bc5842d27725b7fbbe6b5006) controller: update registry API calls for 0.10.x+
 - [`042ef60`](https://github.com/deis/deis/commit/042ef60e12bf00a1d5596e725dcf9ca14053e068) builder: properly escape backticks in envvars
 - [`93bb0fd`](https://github.com/deis/deis/commit/93bb0fd9cb33e5b8bdcfdc277d15d61b938a88d4) router: disable SSLv3 CVE-2014-3566
 - [`7c4fc31`](https://github.com/deis/deis/commit/7c4fc31dc8565b7f992ac5121f40eecb63193c1a) go: go 1.4 static binaries
 - [`345f900`](https://github.com/deis/deis/commit/345f900a636012ea8461626ae9ec9a676254106a) deisctl: add custom error message for global units.
 - [`37cc14f`](https://github.com/deis/deis/commit/37cc14f77c1cf64e3482b7507662b5156b180b20) user-data: ensure nf_conntrack kernel module is loaded
 - [`caa3eac`](https://github.com/deis/deis/commit/caa3eac70d85c632bf4dab7ec9ee773007603301) controller: ensure cleanup of "deis run" fleet units
 - [`05c322c`](https://github.com/deis/deis/commit/05c322c675c7ff6007a40f3d7de6554489f35281) controller: kill processes removed from procfile
 - [`bad8a73`](https://github.com/deis/deis/commit/bad8a731e3379d1eb9d336a21dc9b50dece1f00e) docs+controller: update botbot.me and stackoverflow links
 - [`0a206fd`](https://github.com/deis/deis/commit/0a206fdb6902006efe54c7e1e83b9a32de8288e7) tests: source /etc/environment only if it exists
 - [`96c7a4c`](https://github.com/deis/deis/commit/96c7a4c555e05d2eb92e82a20babe6234342415c) publisher: remove unnecessary BUILD_IMAGE variable from Makefile
 - [`458be25`](https://github.com/deis/deis/commit/458be251c3d9918dd0725adfa33881806c80c36b) controller: ignore run procs when updating structure field
 - [`1eb8fac`](https://github.com/deis/deis/commit/1eb8fac66fffe6221019945585da25c1f5632e40) contrib/azure: bump load balancer timeout to 20 minutes
 - [`a49ff60`](https://github.com/deis/deis/commit/a49ff604d73a248710b849679e6531f3975ee5c6) tests: check for expected prompt
 - [`0901332`](https://github.com/deis/deis/commit/0901332cad68adfa4f6ca92b55927aa13642cdff) contrib/ec2: restrict EBS type to standard or gp2
 - [`4295eea`](https://github.com/deis/deis/commit/4295eeaa8effe733fe6955ec294fd428972bc239) client: display default answer for auth:cancel prompt
 - [`93ea9aa`](https://github.com/deis/deis/commit/93ea9aa8577670e4bd2d9144393cde511e139178) controller: remove timeout on container launch
 - [`9793fd5`](https://github.com/deis/deis/commit/9793fd506e519e8a86ebcbcf106ae7cff79a87b3) contrib/azure: fix health-check for builder
 - [`a820841`](https://github.com/deis/deis/commit/a820841aa5ddd27b80ce3de7896e638ad94c4115) builder: ignore header when removing dangling docker images
 - [`b590f4c`](https://github.com/deis/deis/commit/b590f4c190ebd680d6ea6796771125260d31457c) builder: ensure loopback devices for docker's devicemapper
 - [`2ace8be`](https://github.com/deis/deis/commit/2ace8be56d9a5cdacc47cd84796f0ae883a41008) builder: use Docker's default storage driver
 - [`c4efc67`](https://github.com/deis/deis/commit/c4efc67f1e79c881c2a6144942ed2975c29a86a8) publisher: use godeps
 - [`9090c77`](https://github.com/deis/deis/commit/9090c7770b95f0b230e245afd1be6651c0880d6d) logger: create logRoot on startup
 - [`555f863`](https://github.com/deis/deis/commit/555f8639cddbac9fbab97e4390bbf1eb97e005f4) cache: use godep to build binaries
 - [`1610707`](https://github.com/deis/deis/commit/1610707e6df5ee58125d616316ed066d8c6d0a6f) builder: properly escape backticks
 - [`76c5b0e`](https://github.com/deis/deis/commit/76c5b0ef9e782480f9a528214daf9f93151ed04f) (all): force tags in Makefiles
 - [`a002dd4`](https://github.com/deis/deis/commit/a002dd44024801610492ecd5bd2966e8d4719311) controller: improve domain name validation per RFC 1123
 - [`1d4eaf5`](https://github.com/deis/deis/commit/1d4eaf53956841a1fa7ffc9ec2a746b351e6c514) builder: Fixed invalid redirect syntax

#### Documentation

 - [`406075b`](https://github.com/deis/deis/commit/406075b5c25667fad2f8616e29f156b01ac727a8) contrib: add link to new community project
 - [`17752b5`](https://github.com/deis/deis/commit/17752b5e4c8c29e56ee2cb92e23f00d98c2db9c6) contributing: remove reference to `deisctl ssh`
 - [`440bb16`](https://github.com/deis/deis/commit/440bb1676bb978c0d07d93b622fc1c3c731f57f3) contrib/README: add link to New Relic unit for CoreOS
 - [`33ce22d`](https://github.com/deis/deis/commit/33ce22d599e926a6e5e95c0b30ca7cd9a4625866) contrib: link to community projects
 - [`5f6ce0b`](https://github.com/deis/deis/commit/5f6ce0baba25eea0bd5e51b33bf0fb60059b8575) managing_deis: add proxy docs
 - [`41fbc85`](https://github.com/deis/deis/commit/41fbc85f329db6a6ab9a17ee126a32f8e9a553bf) installing_deis: link to Azure on quick start guide
 - [`1d26154`](https://github.com/deis/deis/commit/1d2615418530c710baef72b96222ac0def0590c8) managing_deis: fix typo in Recovering Ceph Quorum
 - [`59bb22d`](https://github.com/deis/deis/commit/59bb22dd1af49ca7b58f9c6a4ca14e0bc2296af6) (all): add disk usage docs
 - [`aeeb430`](https://github.com/deis/deis/commit/aeeb430319448f4e9073fc8d703a11f9761b568f) managing_deis: use non-forked s3cmd
 - [`cede565`](https://github.com/deis/deis/commit/cede565cd8a31543a3e735e06fdd9947ee0b3d43) logspout: Fix remote syslog example

#### Maintenance

 - [`d97f502`](https://github.com/deis/deis/commit/d97f50263810aa3983406e48a4ecf814fc77cd8a) release: update version to v1.4.0
 - [`b5335b4`](https://github.com/deis/deis/commit/b5335b45cfabd65fbbb841b6bb985c93d36167ca) tests: update CI node setup instructions
 - [`4d3c4b4`](https://github.com/deis/deis/commit/4d3c4b47d7f6685de8e507e91609e0331c79a231) (all): bump to CoreOS 557.2.0
 - [`2f45df9`](https://github.com/deis/deis/commit/2f45df99cf70bb9b60ddd146f1cf23284d68aae0) builder: update Docker to 1.4.1
 - [`a5661da`](https://github.com/deis/deis/commit/a5661da7a68d8cd23cf809a4270a639dba33a957) registry: update docker-registry to v0.9.1
 - [`05ad53c`](https://github.com/deis/deis/commit/05ad53c085b90c8eb98c9be02bc7308cedaa40f0) controller: update docker-py to 0.7.2
 - [`e90633a`](https://github.com/deis/deis/commit/e90633aa71c21615f4b5f9f69b90ee3ae9998861) (all): update pip installer tool to 6.0.8
 - [`3875bdc`](https://github.com/deis/deis/commit/3875bdcd2a4f643f141ef71ddb1dfb648c04f000) registry: bump docker-registry to 258398d
 - [`3aed97b`](https://github.com/deis/deis/commit/3aed97b2e48b082465e215d3e873952c4d6ea9eb) controller: update djangorestframework to 3.0.5
 - [`4c02b5b`](https://github.com/deis/deis/commit/4c02b5b27779bde062e82806e8016006690d622e) contrib/azure: bump to CoreOS 522.6.0
 - [`1830431`](https://github.com/deis/deis/commit/18304318bbaaa97a76bc7f4cae229829f62ad5f0) contrib/azure: bump Docker volume to 100GB to match AWS
 - [`51dcbe2`](https://github.com/deis/deis/commit/51dcbe2211b9402658326d51c7a15cec6edcbeaf) controller: update PostgreSQL driver to 2.6
 - [`fca5610`](https://github.com/deis/deis/commit/fca5610b90ac8334becb9918fe8193c746b7a2b1) client+controller: update flake8 code checker to 2.3.0
 - [`6c2c7e8`](https://github.com/deis/deis/commit/6c2c7e8e24363f11875ee7a9005aaf7ff318ce1f) deisctl: update vendored fleet code to 0.9.0
 - [`f540e96`](https://github.com/deis/deis/commit/f540e96bf5b6c2c4f8de6e0c4c7ca3bbfd50a1cc) controller: No sudo for deis
 - [`65fbebb`](https://github.com/deis/deis/commit/65fbebbcc18dd353b6e89a2deb0c506cc2f86f3d) (all): update to go 1.4.1
 - [`6d313a6`](https://github.com/deis/deis/commit/6d313a60259ed41247ab67a4f534c776a2f5778b) client: remove deprecated settings file converter
 - [`bf92c4c`](https://github.com/deis/deis/commit/bf92c4c7b882394a421fc5858b3388916f97bb6a) logspout: remove unused Dockerfile
 - [`dab12f4`](https://github.com/deis/deis/commit/dab12f45f5f7b16cd794a57cd8630f3754a79e9c) controller: update gunicorn to 19.2.1
 - [`5e951c7`](https://github.com/deis/deis/commit/5e951c71cad7edec993aec6fa4ad8a1263718d0b) release: update version in master to v1.3.1
 - [`19125af`](https://github.com/deis/deis/commit/19125af9430c42b58b5e8e149a5a072a2dd34d5a) logspout: remove checks from make clean
 - [`7480d70`](https://github.com/deis/deis/commit/7480d708da565b245dba1ab7993ec2e359b63b0c) builder, logger, logspout: remove unused references to BUILD_IMAGE from the Makefile
 - [`9b77f0f`](https://github.com/deis/deis/commit/9b77f0fc689cc8603526b6d29b57d59ca822c4a6) release: update version in master to v1.4.0-dev
 - [`a4081f4`](https://github.com/deis/deis/commit/a4081f45f2a0ef4e2530a685d9d509a6d5638340) Vagrantfile: sync with changes in coreos-vagrant
 - [`3997d1b`](https://github.com/deis/deis/commit/3997d1b4a24a7269f189e12b971c94c22dbf90cf) controller: update djangorestframework to 3.0.4
 - [`1fb501b`](https://github.com/deis/deis/commit/1fb501b9569a18d4fecd4c0336d7eab5bdca2bde) builder: Updated gradle and play buildpack

### v1.3.0 -> v1.3.1

#### Fixes

 - [`71bd706`](https://github.com/deis/deis/commit/71bd7060836e853efac41b4e9a990daddc487ccf) cache: use godep to build binaries

#### Maintenance

 - [`deaa34c`](https://github.com/deis/deis/commit/deaa34cd0d83ddb7107d4862dfb3475f49b4672a) release: update version to v1.3.1

### v1.2.2 -> v1.3.0

#### Features

 - [`dc5b866`](https://github.com/deis/deis/commit/dc5b8664a77f0a0947f7c971c091e3c845431517) deisctl: support store-admin in refresh-units
 - [`a238964`](https://github.com/deis/deis/commit/a23896420a7ad6b18c50e70cdefa2f01966033ff) store: add store-admin container
 - [`678f20b`](https://github.com/deis/deis/commit/678f20be1aa09635db72ca8406e006244c59d8fa) client: add flag to disable SSL certificate verification
 - [`bafc3e2`](https://github.com/deis/deis/commit/bafc3e2e956ede7bc9597c8dc1ef1a5161865f62) router: listen on port 80 for platform SSL
 - [`6117570`](https://github.com/deis/deis/commit/6117570a9842a9eaf829ba083dadad2a82fe43b7) builder: disable simultaneous push in the same repository
 - [`a30137e`](https://github.com/deis/deis/commit/a30137e016ad1ef2da2fd07e808347e93040cb79) (all): support RFC 6598 for insecure registries
 - [`b8ff429`](https://github.com/deis/deis/commit/b8ff429c9457395478f79b0c994da5d17bcb5def) client: Add controller to whoami command.
 - [`468cb1f`](https://github.com/deis/deis/commit/468cb1f201b35452feaaf33206900764d6086e59) slugbuilder: update nodejs buildpack to v64 (yoga release
 - [`8477f3e`](https://github.com/deis/deis/commit/8477f3ed8b74c166e9a1dd8a02fe077b609b77f9) store: add gateway mock
 - [`4b791d2`](https://github.com/deis/deis/commit/4b791d29fd9f3e7ea78e801ae409355c73018ea3) controller: add hostname exposing to units
 - [`72fac56`](https://github.com/deis/deis/commit/72fac56a16797caeca70054d9c1d6882fcd19fe3) contrib/ec2: Set ELB timeout to 1200s
 - [`1ed1281`](https://github.com/deis/deis/commit/1ed1281d71d85330a123feab481c3f9d60448aac) docker: add env DEIS_RELEASE

#### Fixes

 - [`c20d892`](https://github.com/deis/deis/commit/c20d8925decbcfb9db5a99c1c92de3527312c9be) logger: add logger build to gitignore
 - [`dc0dae9`](https://github.com/deis/deis/commit/dc0dae9320434050665ed85b12315eb2b038016c) contrib: use 1GB per VM"
 - [`1a8a656`](https://github.com/deis/deis/commit/1a8a656b0c111883835b28ef164d71f04592cc82) logger: use godep for dependencies
 - [`d39c9c8`](https://github.com/deis/deis/commit/d39c9c85db3a01d9b949edd991516536a2e3a706) docs: refer to correct "deis perms:create" command
 - [`90fe432`](https://github.com/deis/deis/commit/90fe432dcf5e436802f0cb4bbf2c2492dd082e30) contrib/coreos: etdcd_request_timeout value type
 - [`c67afb4`](https://github.com/deis/deis/commit/c67afb4ea9f65e7e6b3130aba151f11b23c520ae) builder: correctly handle error on request timeout
 - [`7af2e16`](https://github.com/deis/deis/commit/7af2e16b8d3def0b0477f7c8f24b8c18e8473594) security: increase max conntrack connections
 - [`aabb08f`](https://github.com/deis/deis/commit/aabb08f2d807a8a5b94283b5d45d8e4e9487ec3e) controller: return proper json messages on error
 - [`4692bda`](https://github.com/deis/deis/commit/4692bda95ba59fcc5034b3897db3f4f5343792c0) controller: validate key material on upload
 - [`861108a`](https://github.com/deis/deis/commit/861108a1fe5fdbd2bcc9576a834d3a7083c53f39) controller: return 403 when a user does not have permission
 - [`43d8a92`](https://github.com/deis/deis/commit/43d8a9207e9493625dafc92e2c95ee36dcec0bf3) builder: exit if an error occur compiling go
 - [`283ead3`](https://github.com/deis/deis/commit/283ead314f2a4d7bd79560361418b9d189d9e129) client: check for expected destroy response
 - [`31052eb`](https://github.com/deis/deis/commit/31052eb18ea780d2ada7ef09fbaa7a30d70c4f1f) database: increases max_connections to 250
 - [`2515661`](https://github.com/deis/deis/commit/25156613c820d2ed753523c9ec44bbd8b8d72410) store: use hostnames for mon host
 - [`f2d2f2f`](https://github.com/deis/deis/commit/f2d2f2f6b0eae1ea3195c0ba6591391dcb50c486) builder: allow procfile to override default process types
 - [`7b8e6d7`](https://github.com/deis/deis/commit/7b8e6d7acd734d8a9be4f63ec53ce1fe4aedd6c2) controller: track create and start separately to catch early errors
 - [`a23bd6a`](https://github.com/deis/deis/commit/a23bd6a6893885d64094458eec60398ea2f34918) database: fix build because of python-daemon
 - [`70af19c`](https://github.com/deis/deis/commit/70af19c47d7452fbef51bcca6eb2534f671c4d23) (all): bump version in Dockerfiles
 - [`c85b062`](https://github.com/deis/deis/commit/c85b06219cbaa4876b9321acc7007ac2a1bde05f) client: remove fetch call from apps:create

#### Documentation

 - [`db7a930`](https://github.com/deis/deis/commit/db7a930967f65c8c25d0f844e4cb064b87f22cbe) ec2: Update installing_deis/aws
 - [`88a2de2`](https://github.com/deis/deis/commit/88a2de28ed62f5dbbfcda99ea976cccfb0858305) ec2: Document update_ec2_cluster.sh
 - [`249cded`](https://github.com/deis/deis/commit/249cded318147a54ea4327a71a27d236d5ccfd5e) installing_deis: add --use_compute_key option flag
 - [`102633e`](https://github.com/deis/deis/commit/102633eb1f94130d859850b7a77f684a0968bcfa) contributing: customize test-acceptance on component addition
 - [`4ef7b4d`](https://github.com/deis/deis/commit/4ef7b4d269d5568356df476122b3183c9d995a34) faq: added a link to unofficial Chinese documentation
 - [`d7d3f6d`](https://github.com/deis/deis/commit/d7d3f6dc7cb7e868243aa631fa3fac01272d018b) (all): add Ceph quorum documentation
 - [`cfae0e3`](https://github.com/deis/deis/commit/cfae0e31ac9dd8587b4f6cdaf94895e0f6800789) (all): reference deis-store-admin for admin tasks
 - [`aba69b1`](https://github.com/deis/deis/commit/aba69b14d12ed012eadd37d508398b8995f7f326) system-requirements: Fix broken links
 - [`42931ab`](https://github.com/deis/deis/commit/42931abccb7ece448f19ed49521687033b896f3c) ssl: add instructions for configuring ssl on routers
 - [`3998c15`](https://github.com/deis/deis/commit/3998c156485aa20a44616ac048d638daac04515e) private-networks: add RFC 6598 address range
 - [`0cbd091`](https://github.com/deis/deis/commit/0cbd0911d2fd1ba1338a9895e4d93ca2c1e80d35) managing_deis: clarify reasons for disabling CoreOS upgrades
 - [`ec07233`](https://github.com/deis/deis/commit/ec07233fa27b740dad113f3e508312cd349cf5b1) using_deis: add disabling registration docs to user creation
 - [`6022e6c`](https://github.com/deis/deis/commit/6022e6c994fd2bf82a0d395e59fc78096ccb69bf) controller: add unit hostname documentation

#### Maintenance

 - [`dddfec3`](https://github.com/deis/deis/commit/dddfec39ace1016e5fa9ef9a82fda05f61415704) release: update version to v1.3.0
 - [`a14a4ed`](https://github.com/deis/deis/commit/a14a4ed8beaf94ade91e2b1605b11b72680110fd) contrib/rackspace: bump CoreOS to 522.6.0
 - [`e4fc01c`](https://github.com/deis/deis/commit/e4fc01c76cd878517e273c65bcae83b446904c2d) (all): bump CoreOS to 522.6.0
 - [`863acd9`](https://github.com/deis/deis/commit/863acd903542fc1970011890e4a559f99e0c8f41) builder: bump heroku-buildpack-php to v57
 - [`60a61cd`](https://github.com/deis/deis/commit/60a61cde9ec1cd6710aa56158556928bfad19163) contrib/rackspace: bump CoreOS to 522.5.0
 - [`2a56376`](https://github.com/deis/deis/commit/2a5637669cc31cd036955a7b596a88703414ae1b) controller: bump API version to 1.1.2
 - [`73a9c08`](https://github.com/deis/deis/commit/73a9c08aeff439d1904108348010dbf17915ad12) controller: update static wsgi library to v1.1.1
 - [`b174fae`](https://github.com/deis/deis/commit/b174faeeb72f1c7cfd42eaf2545474573b07f104) controller: bump django-rest-framework to 3.0.3
 - [`ffac129`](https://github.com/deis/deis/commit/ffac1293fa5bed2faaeeccdafaccd511202704be) release: update version in master to v1.2.2
 - [`60b1aa3`](https://github.com/deis/deis/commit/60b1aa31388fadc59353ce34542581e152600783) controller: update django-guardian to 1.2.5
 - [`898193a`](https://github.com/deis/deis/commit/898193adedb3645b52de43d8ff4e339b82060866) client: update python-dateutil to 2.4.0
 - [`f0b9358`](https://github.com/deis/deis/commit/f0b9358cd1e2bc4c9ce1b3f6502532b8ddb730a4) controller: update paramiko to 1.15.2
 - [`544f55b`](https://github.com/deis/deis/commit/544f55bf96300a398fe8d15d7a5e5cbb1195210f) (all): bump CoreOS to 522.5.0
 - [`dbd42a7`](https://github.com/deis/deis/commit/dbd42a7352759f9ac98cbd59bf20ac757672bc6d) release: update version in master to v1.2.1
 - [`a57cabf`](https://github.com/deis/deis/commit/a57cabfea0a894ff4e5909301b72e8883d283dd2) controller: update Django to 1.6.10 security release
 - [`0db7905`](https://github.com/deis/deis/commit/0db79055c99532aa7bd43c23bfadc3951b883240) release: update version in master to v1.3.0-dev
 - [`f67dd7f`](https://github.com/deis/deis/commit/f67dd7fe621060b27e59984f66f7140a9601193e) controller: update Django to 1.6.9 bugfix release

### v1.2.1 -> v1.2.2

#### Fixes

 - [`97d9406`](https://github.com/deis/deis/commit/97d9406f397cbf794b39a133a408fb0aa2900731) database: fix build because of python-daemon
 - [`5d508c6`](https://github.com/deis/deis/commit/5d508c6882e6767f8d9c2982d004881b1d644185) builder: allow procfile to override default process types

#### Maintenance

 - [`7bb8caa`](https://github.com/deis/deis/commit/7bb8caac471dfd58fa3533e0d099f7e86b698776) release: update version to v1.2.2

### v1.2.0 -> v1.2.1

#### Maintenance

 - [`3292b7c`](https://github.com/deis/deis/commit/3292b7c693b84f5789bbb12bee33ad61f2ec400d) release: update version to v1.2.1
 - [`0d13d66`](https://github.com/deis/deis/commit/0d13d66ee40b6d41d6328b411b98521283a1300b) controller: update Django to 1.6.10 security release

### v1.1.1 -> v1.2.0

#### Features

 - [`73f2eda`](https://github.com/deis/deis/commit/73f2eda31c758fab0893b24a1cda92134dd352ac) contrib/ec2: use ephemeral drive for etcd journal
 - [`6db1e50`](https://github.com/deis/deis/commit/6db1e50edef6cc7802988250643f9649955bbadb) contrib/ec2: use absolute path names for scripts
 - [`41d1cc9`](https://github.com/deis/deis/commit/41d1cc9b9a81528a74cfdd6f173826224abc554d) contrib/ec2: use separate volume for /var/lib/docker
 - [`275ef40`](https://github.com/deis/deis/commit/275ef40d15baf4554521055b41fe6857dbcaec92) pkg: add pkg time
 - [`331e4b4`](https://github.com/deis/deis/commit/331e4b466c68c8dd53bce7fc341946bb0a4f0425) contrib/coreos: add etcd debug unit
 - [`16003d9`](https://github.com/deis/deis/commit/16003d9736669aa46359b4151b3c11743f1dbb43) builder: add cedarish/build target
 - [`c360729`](https://github.com/deis/deis/commit/c3607298c3a837160a9d9c4e58f800b4fee09da8) Makefile: add dev-cluster command
 - [`e8b61f8`](https://github.com/deis/deis/commit/e8b61f89b4ea4d8a122aa3229c50c9a9fdfc0d3f) registry: deis-cache for saving meta data
 - [`828bca9`](https://github.com/deis/deis/commit/828bca9fa88959769e92ded8637658a3114ecddb) client: add "--buildpack" option to "deis create"
 - [`f6d121e`](https://github.com/deis/deis/commit/f6d121e202ecff96c89807fb1deec4e5d8e2caf4) vagrant: check for discovery url
 - [`997bd50`](https://github.com/deis/deis/commit/997bd50b8a2fc8446cf32176e5571e12d56b9d69) builder: do not lookup the remote host in ssh connections
 - [`17f9ee6`](https://github.com/deis/deis/commit/17f9ee66ff3e06ea10e400d777e85b2646e6282a) security: allow the access to new nodes if we are using the custom firewall
 - [`28828a7`](https://github.com/deis/deis/commit/28828a7c3f216a46e8e39587324d300a349169a5) contrib/ec2: add optional support for IAM instance profiles

#### Fixes

 - [`a6b7d86`](https://github.com/deis/deis/commit/a6b7d86903f3ed0f207d6d600dfdfe4ccfdc80df) contrib/ec2: prepare etcd before starting etcd
 - [`7fa8358`](https://github.com/deis/deis/commit/7fa83584227d3df76cb8128e14aee5af503849fe) controller: update tests to check for 403
 - [`a130937`](https://github.com/deis/deis/commit/a130937c45f37727c559e61a2c189a33ff78ba03) controller: disallow unauthorized users from rolling back releases
 - [`33f6e6a`](https://github.com/deis/deis/commit/33f6e6af551e12067d203228cc346ed896cdd1af) controller: disallow unauthorized users from deleting apps
 - [`a9cc25c`](https://github.com/deis/deis/commit/a9cc25cd670b2fecd33cbd761fe09fedc4a4404a) contrib/coreos: fix debug-etcd command
 - [`70e92f2`](https://github.com/deis/deis/commit/70e92f26c8556b1c17fac92c0dd4087f9f2876c4) registry: use empty strings and not nil in config
 - [`14353e1`](https://github.com/deis/deis/commit/14353e1cdad518526a4cee3e1b726615a84e2070) controller: disallow unauthorized users from modifying apps
 - [`42966f5`](https://github.com/deis/deis/commit/42966f528dc369f03d2fd9227fd20d18cc384fc9) controller: adds migration to get rid of dockerfile's u'False' values (upgrade fix
 - [`ba5a3aa`](https://github.com/deis/deis/commit/ba5a3aa9a1d846bf7494c2456844f2b9cf71f3c9) controller: disallow evil users from creating builds
 - [`d3526aa`](https://github.com/deis/deis/commit/d3526aad495d03371b478f74238c744431d457db) client: encode unicode properly for apps:run
 - [`4c1bd3f`](https://github.com/deis/deis/commit/4c1bd3f31a78e3bb91a3304f8966d5dcd820bf97) .gitignore: hide local vagrant config.rb from version control
 - [`d4e45ce`](https://github.com/deis/deis/commit/d4e45ceae6b06358d661d63fbcb9425349a12912) slugbuilder: allow users to cat slugfile
 - [`5395fe7`](https://github.com/deis/deis/commit/5395fe71a25c1e18454900778a954439194decd8) router: fix deploy Makefile target
 - [`37b5f78`](https://github.com/deis/deis/commit/37b5f78ec8b1c1f28707f40ca266d374d2858f34) setup-node.sh: create database role for controller unit tests
 - [`3a6f145`](https://github.com/deis/deis/commit/3a6f14561384aa6f93d803adc286bf0d940720a1) router: missing check file in git ignored by .gitignore
 - [`2ca5f35`](https://github.com/deis/deis/commit/2ca5f353a6a4f9588d300c7d48e277ec73a0f591) deisctl: unit test should also use main version package
 - [`a997083`](https://github.com/deis/deis/commit/a9970836514ed742eb743ffb7d4912e5127094f5) deisctl: fixup most of `go vet` errors
 - [`1a0dc0b`](https://github.com/deis/deis/commit/1a0dc0bf89b1d4518db374e2ba8bf88c8319b200) deisctl: fixup unit tests
 - [`fbc6b39`](https://github.com/deis/deis/commit/fbc6b3989eb10cba7d687e9fad9f9f6b754c6869) deisctl: hook up unit tests
 - [`b1fff20`](https://github.com/deis/deis/commit/b1fff20872bd33b9058ad2e7c39a5dce130e6274) contrib: use tab indents for Makefiles
 - [`dff278c`](https://github.com/deis/deis/commit/dff278cdc0e3f240328b5a24e8d0381bfd979c12) Makefile: remove CLIENTS from release task
 - [`f5a8ffb`](https://github.com/deis/deis/commit/f5a8ffbaed3384eadaea404ba1065fd5d31e7653) deisctl: properly resolve ~ in path
 - [`026ebbc`](https://github.com/deis/deis/commit/026ebbc1b8c5a8aab051a5d7921054c9bf3588b9) builder: delete images on destroy
 - [`d09e3f6`](https://github.com/deis/deis/commit/d09e3f62387de6e1548e5d884fbde2875a7ea05a) deisctl: remove unimplemented --verbose flag from "deisctl config"
 - [`d9de0e4`](https://github.com/deis/deis/commit/d9de0e4ebd284da5ece63110cebc1fc543c3803f) deisctl: parse global and command-specific options correctly
 - [`10f80cf`](https://github.com/deis/deis/commit/10f80cf8528cb669780bf686bafdd7108ec5e83a) contrib: use 1GB per VM
 - [`a55c123`](https://github.com/deis/deis/commit/a55c12397981fd3d96c21a5af2a935230aa6fa4f) Vagrantfile: disable replacing of default insecure key
 - [`695e53e`](https://github.com/deis/deis/commit/695e53e168088c3c45e5c348adeaff9c444e3378) tests: destroy without --app flag to remove git remote
 - [`2358fcd`](https://github.com/deis/deis/commit/2358fcd2f497166312783c55aa6f5029f0d541c8) store-gateway: don't log to stdout
 - [`690b380`](https://github.com/deis/deis/commit/690b3807763e40cc6f93d30763a5f532c8db4cd0) store: list all monitors in [global] section of config
 - [`932a251`](https://github.com/deis/deis/commit/932a2519a51a2823efc55a76dba3b41961793ce9) client: update PyInstaller to package requests lib properly
 - [`6a2be87`](https://github.com/deis/deis/commit/6a2be872dc0a595797d65fa657ac49ba74e1a2f9) controller: set config as the issuing user
 - [`239f79d`](https://github.com/deis/deis/commit/239f79d677c0f7246626b73ab4b7868e0dc818f6) builder: the value could be a number
 - [`56e4cec`](https://github.com/deis/deis/commit/56e4cec4763e96843af611d47373193091d5a6e0) gateway: wait until radosgw is accessible
 - [`46a2159`](https://github.com/deis/deis/commit/46a2159aade14e8dde80b9c520d6b8c0a2b0c1b2) deisctl: prevent build on windows
 - [`68ab344`](https://github.com/deis/deis/commit/68ab344bfce754b08a31798c05da83ffcc13a4da) client: avoid git remote deletion for apps specified by -a

#### Documentation

 - [`0df3128`](https://github.com/deis/deis/commit/0df3128dcb041d29a7688b68906553645b94b878) releases: update with major/minor and patch procedures
 - [`eb113c9`](https://github.com/deis/deis/commit/eb113c978f1efd8986845df39c7cfd5498807399) using_buildpacks: add sidenote on Scala buildpack
 - [`6799be0`](https://github.com/deis/deis/commit/6799be025cb90f6f12bde1bac75e5e5076fc88ab) customizing_deis: remove slugbuilder/slugrunner settings
 - [`8b582d4`](https://github.com/deis/deis/commit/8b582d4c42c827e272bcab1ef82108a35437a6a4) docs/managing_deis/platform_logging: change example to prefer tcp/tls
 - [`08602f1`](https://github.com/deis/deis/commit/08602f1a29522daf3494c6cd99a64f42bd2066cf) troubleshooting_deis: document debug-etcd service
 - [`38c9379`](https://github.com/deis/deis/commit/38c937933a1a9d8daad28038a9f35d97021304bc) disable-registration: adds documentation on how to disable user registration
 - [`240e0b7`](https://github.com/deis/deis/commit/240e0b7922f01914f8a55de6919c6c7d65573f7e) openstack: move OpenStack docs to docs site
 - [`0351b99`](https://github.com/deis/deis/commit/0351b993bddec6af9331fff23681a83729707938) openstack: fix instructions
 - [`7ba347f`](https://github.com/deis/deis/commit/7ba347f3cb4dddb08e9e59460b7e581aef742308) hacking+README: clarify that deis CLI requires python 2.7
 - [`b471c36`](https://github.com/deis/deis/commit/b471c3653b617a4483398a39709dab7525312eed) (all): use https:// URL for ci.deis.io
 - [`3b772fa`](https://github.com/deis/deis/commit/3b772fa00b4e2e2006b67511c5b18a6e78267b20) contrib: Clarify source of DEV_REGISTRY IP Address
 - [`344f892`](https://github.com/deis/deis/commit/344f892030d85c2a32162a3176b415ffab55b058) using_deis/(all): add scale example and clarify process == container

#### Maintenance

 - [`170effe`](https://github.com/deis/deis/commit/170effe7045f377084f51f370a7df09efced565b) release: update version to v1.2.0
 - [`49ed8f8`](https://github.com/deis/deis/commit/49ed8f8bfb4ed629768d5ef77b14d90a415f2cf0) controller: Fix typo on README.md
 - [`64d735d`](https://github.com/deis/deis/commit/64d735d3ba4d3368dad0cb8e49d9f69dc22fe007) registry: bump registry to eb62607
 - [`2e94e9d`](https://github.com/deis/deis/commit/2e94e9d2f4d4d028fe405d896317e599feb5fd68) cache: update Redis to 2.8.19
 - [`d3efd3d`](https://github.com/deis/deis/commit/d3efd3d10fe19611b2cb497c828a22e61bba8b89) controller: update South to 1.0.2 bugfix release
 - [`9968c23`](https://github.com/deis/deis/commit/9968c2356d8feb84393af2fd95e0015da9aab788) contrib/rackspace: bump to CoreOS 494.5.0
 - [`1c6e55f`](https://github.com/deis/deis/commit/1c6e55fe9f48c2a172b597566f29b6a2fd8d7e14) slugbuilder: bump heroku buildpacks
 - [`c706fb4`](https://github.com/deis/deis/commit/c706fb4eb326b4561bb8bc2b3e17179571833f7d) client+controller: update python requests lib to 2.5.1
 - [`6a1ec31`](https://github.com/deis/deis/commit/6a1ec31c74015604070f1cf42b9c4e8cd6f6e044) controller: update django-cors-headers to 1.0.0
 - [`8f8a226`](https://github.com/deis/deis/commit/8f8a226267a13e0af7a20ec2ada89bd594797e09) builder: bump cedarish
 - [`799c32b`](https://github.com/deis/deis/commit/799c32b32cb35dea89486b7ef6dedbebaab62846) store: remove redundant Makefile task
 - [`65b39eb`](https://github.com/deis/deis/commit/65b39ebb5cfa8490d5541407d4b20341597c1efc) release: update release version in master to v1.1.1
 - [`8174622`](https://github.com/deis/deis/commit/8174622bb3db9e45a983f0e0bb902360c3b3cb79) (all): update base CoreOS to 494.5.0
 - [`5ebf3f2`](https://github.com/deis/deis/commit/5ebf3f2dcd5009ed49c82538237801ed81cdb5de) (all): Allow testing against CoreOS alpha
 - [`be9d842`](https://github.com/deis/deis/commit/be9d842605410ffe97812d997da2086a7f04c595) builder: update Docker to 1.3.3
 - [`f21159f`](https://github.com/deis/deis/commit/f21159fc1a9c440f8a29e1862b4a1764f4dfcea2) release: update version in master to 1.2.0-dev

### v1.1.0 -> v1.1.1

#### Fixes

 - [`fe56989`](https://github.com/deis/deis/commit/fe56989b2a2df5143a81137bffc751b3728e389c) deisctl: remove unimplemented --verbose flag from "deisctl config"
 - [`430edbe`](https://github.com/deis/deis/commit/430edbeebea436c2d1ef726a479b92fa9020cf78) deisctl: parse global and command-specific options correctly
 - [`962a181`](https://github.com/deis/deis/commit/962a1814c6aed02ce8c1166126e244918efb10ad) builder: the value could be a number

#### Maintenance

 - [`c9c181e`](https://github.com/deis/deis/commit/c9c181e9b5a1cd58a5fbab70446ea43fd432a0a1) release: update version to v1.1.1
 - [`6fcb278`](https://github.com/deis/deis/commit/6fcb2789a93361c59ebb446d69941e7ec43cd954) (all): update base CoreOS to 494.5.0
 - [`3075405`](https://github.com/deis/deis/commit/3075405d66d347c4e132f5ac48897a43b7ef78eb) builder: update Docker to 1.3.3

### v1.0.2 -> v1.1.0

#### Features

 - [`98c8e2e`](https://github.com/deis/deis/commit/98c8e2ed9d36fe3cf7eaf89a0ff05b9b63ae46a6) database: disable copy on write to improve performance"
 - [`d6246f9`](https://github.com/deis/deis/commit/d6246f93dc3e54112605f0aa1271d3b23699a003) controller: introduce dockerignore
 - [`1777203`](https://github.com/deis/deis/commit/1777203299262c4102a0a38969cd731ccae51cd2) (all): provision using CoreOS stable channel
 - [`3c52611`](https://github.com/deis/deis/commit/3c5261129b024ebe56dfa3c9a6c345ff94a42214) contrib: add PREFIX option to digitalocean script
 - [`1cbb524`](https://github.com/deis/deis/commit/1cbb52474e23b043ce50ad41818649102f8e7db1) deisctl: add -ssh-timeout command-line option
 - [`38e8e29`](https://github.com/deis/deis/commit/38e8e290a0686df778bbaab0001f95543fc658d0) tests: upload dumped logs to S3
 - [`f729303`](https://github.com/deis/deis/commit/f7293030df04903d13b4269d1b6db4772cfeacbb) (all): provision using CoreOS beta channel
 - [`4a1c1f9`](https://github.com/deis/deis/commit/4a1c1f99eb9cb881702bf57cbcfb39389d0c1b20) tests: dump all component and app unit logs on failure
 - [`1017388`](https://github.com/deis/deis/commit/10173885126cf1850d5372e3f7e0260e8644359f) router: include a firewall to mitigate security problems.
 - [`fa44947`](https://github.com/deis/deis/commit/fa44947402d652f433144fa2e5db0eb0f4081ef0) tests: provide debug output on etcd/fleet timeout
 - [`f4304e3`](https://github.com/deis/deis/commit/f4304e30e742fe63e88dfbd52bddbeeb6f4a1a60) client: upload local Procfile on builds:create
 - [`0259313`](https://github.com/deis/deis/commit/025931336e0448fec6a722670dca75ea16dedcee) user-data: enable time synchronization on startup
 - [`fe7f9c2`](https://github.com/deis/deis/commit/fe7f9c2c720a56b7b23d22ea31e62df259f3e932) CONTRIBUTING: add design proposal guidelines
 - [`e1db573`](https://github.com/deis/deis/commit/e1db57333c8d8c31ddf839d164c6d191361b9c29) controller: add X_DEIS_PLATFORM_VERSION header in response
 - [`ce12712`](https://github.com/deis/deis/commit/ce12712f74c041626df4aa6b668b43ff05f6dd5f) controller: check client against API version
 - [`5d6a1c0`](https://github.com/deis/deis/commit/5d6a1c05226f749fb82fc6dae6582eb4c94dbd80) ec2-provision: add script option to set ELB id when running update-ec2-cluster.sh
 - [`1f6f429`](https://github.com/deis/deis/commit/1f6f42980143fb5e39fbaf5d4dd6eea9a4bcf80c) router: support customized server_names_hash_bucket_size setting
 - [`8ce3401`](https://github.com/deis/deis/commit/8ce340131b5a6f544a6045da8f06250757bc2274) client: add pull shortcut to global help
 - [`56ee63c`](https://github.com/deis/deis/commit/56ee63c605df5ccf66ec8890e802f4cad03cb96e) database: disable copy on write to improve performance
 - [`62006c4`](https://github.com/deis/deis/commit/62006c439d5d7ad34486809e8817a5645c5138c8) (all): Scalable registry
 - [`7270aa9`](https://github.com/deis/deis/commit/7270aa9a882f383f318478b0bb2ec2b09bc051fb) router: use godep to include dependencies.
 - [`1e5bd14`](https://github.com/deis/deis/commit/1e5bd140a69ddd7d58fa2a39da6371fc8ab2870d) router: improve log output.

#### Fixes

 - [`b6a2505`](https://github.com/deis/deis/commit/b6a2505490716fa45163af8ae3dbb33a61260541) database: update wal-e""
 - [`d777712`](https://github.com/deis/deis/commit/d77771292b1f8cb6a7d0d4c979c5d42c218ac5fd) database: update wal-e"
 - [`53204c1`](https://github.com/deis/deis/commit/53204c1ec295a1ba968a9e72f9fc628c49a035fa) builder: parse release info correctly
 - [`1e8c17b`](https://github.com/deis/deis/commit/1e8c17b3da252cef550620d3e12da06e89d44212) builder: send nil to controller, not 'False'
 - [`d8b8343`](https://github.com/deis/deis/commit/d8b83435dfddd1a95a57178c8654a61357309d61) controller: revert to `start proctype` for buildpack images
 - [`050263c`](https://github.com/deis/deis/commit/050263c4b86d87a78095345849ea0e9176032364) controller: update App.structure field correctly after scaling
 - [`d9f4e48`](https://github.com/deis/deis/commit/d9f4e4888b855aaab7975f439de40a68bd9a5122) controller: allow users to override `cmd` process types
 - [`043f7e7`](https://github.com/deis/deis/commit/043f7e755798a09522d09d12cfc80b14247dd258) coreos: tune etcd to address high disk i/o environments"
 - [`c04a05f`](https://github.com/deis/deis/commit/c04a05f5d83fd4898bea9ebaca7eba9591a0d5b7) controller: use /runner/init entrypoint only on procfile
 - [`c788826`](https://github.com/deis/deis/commit/c7888268ba60f984131a318305f49600dbda0799) database: update wal-e
 - [`b4d72d2`](https://github.com/deis/deis/commit/b4d72d2d318c5e9dfc3bede1b941c26b23ba0e72) builder: limit cpu shares to 800
 - [`4f95cb8`](https://github.com/deis/deis/commit/4f95cb82c87e75c6492ce0084961d57b0d70ae31) coreos: tune etcd to address high disk i/o environments
 - [`0290f15`](https://github.com/deis/deis/commit/0290f1556175122211eff6e3baddb19f5d1f865e) contrib: use absolute path to user-data
 - [`68f61f1`](https://github.com/deis/deis/commit/68f61f1fee29367fff973143b6e096bb8a75edd1) router: reload nginx after enabling SSL
 - [`d8e86b8`](https://github.com/deis/deis/commit/d8e86b82f4250a74e2a64de320cd1e2c17b6b138) tests: add docs for local testing with s3cmd
 - [`cd37051`](https://github.com/deis/deis/commit/cd370514b8eb0cc179423894a4c4e79cdc5f9b35) contrib/digitalocean: use beta channel on DO
 - [`19cd5f7`](https://github.com/deis/deis/commit/19cd5f76c0d3918e5f46544585c78262816e7412) tests: check for immovable response
 - [`868a7f4`](https://github.com/deis/deis/commit/868a7f43e8ac33ab9329d1b4352ea5ae647a8473) client: log python requests errors
 - [`e568f83`](https://github.com/deis/deis/commit/e568f8341394620f56fd71aab40bbccfc73a01a6) scheduler: send environment to the process
 - [`4183640`](https://github.com/deis/deis/commit/4183640d836eb95b697ddf5b933efa08628b29ba) gce: update to new path for example userdata
 - [`4923487`](https://github.com/deis/deis/commit/4923487b6f2e8d9f1c68e410a2cf32b674fee394) builder: revert to testing local package
 - [`f93fc5f`](https://github.com/deis/deis/commit/f93fc5fcb1ac65d946256ae03fd3486b79fc676c) tests/ps_test: sleep before checking for 503 response
 - [`ef26c5a`](https://github.com/deis/deis/commit/ef26c5a5d0038b66549b3b4c65ee50aa0ec5a9f8) tests: call defer after handling error
 - [`ae01355`](https://github.com/deis/deis/commit/ae01355cb56147af1fb192ca7642ffb1cda09134) builder: switch exit codes to 1
 - [`9871a3d`](https://github.com/deis/deis/commit/9871a3db82b054683cf10fb71cbefbb808be5bf5) docs: remove duplicate PYYaml entries
 - [`981f6cb`](https://github.com/deis/deis/commit/981f6cb3c5491cdf7315d6504d52f29ca7aa5897) controller: read command from procfile
 - [`1d56b0c`](https://github.com/deis/deis/commit/1d56b0c6f4a6b265b6f4e1c567aabd424a225a21) builder: supply empty default process type
 - [`9249b78`](https://github.com/deis/deis/commit/9249b78761106632f0402c82d6a5b1d9f1451d38) builder: avoid issues with non-string release values
 - [`3346d15`](https://github.com/deis/deis/commit/3346d159d531fed2a592c35f6ac145f879dcb3d2) builder: cleanup binaries after `make build`
 - [`1bb19d8`](https://github.com/deis/deis/commit/1bb19d8d76cc7b8a0fe862a00af563138b987988) router/Makefile: use pwd since "docker cp" no longer allows "."
 - [`b5ba7c4`](https://github.com/deis/deis/commit/b5ba7c498ff7dc18b01cdb4746caa080672c10f4) builder: explicitly call out $TMP_DIR on `docker cp`
 - [`b0238d4`](https://github.com/deis/deis/commit/b0238d4c223ead7b2f60f89abf32db5afa36115f) store: don't output secretKey and accessKey
 - [`b7452d6`](https://github.com/deis/deis/commit/b7452d6ebeb456617b7bc98d2d6df3cf3293f3ad) contrib/coreos: fix usage of --insecure-registry
 - [`ae89a2f`](https://github.com/deis/deis/commit/ae89a2f4da9c2fb195853e8ac9f4b60b1e4cfd2b) contrib: add other RFC1918 private addresses
 - [`0868432`](https://github.com/deis/deis/commit/0868432447ae4abc815d6883cc4cff0c9e4f6bab) contrib: add --insecure-registry flag for docker 1.3.2
 - [`8d98ecd`](https://github.com/deis/deis/commit/8d98ecdd17eb85975905ac000bd861b258d76953) (all): Handling of unavailable services
 - [`d668ce8`](https://github.com/deis/deis/commit/d668ce827a0dc63ba8d3493fb564e55d3a7b70a3) client: add "deis shortcuts" to main help message
 - [`35a9434`](https://github.com/deis/deis/commit/35a94340bd565b2374e5ec72e6debd4e6cf7aa55) contrib/util: use /usr/bin/env shebang for bash scripts
 - [`2efcf74`](https://github.com/deis/deis/commit/2efcf743834d1aa332b4000006aed8b38a56dc81) bumpver: reference files which need bumping
 - [`c22550d`](https://github.com/deis/deis/commit/c22550dc2a2d36018fe98746396d8c4573e714d1) bumpver: handle source and docs the same: replace all
 - [`33fb20e`](https://github.com/deis/deis/commit/33fb20e7f94afe164f2ceeb34fc5f2379934f3f7) gateway: Removed unrecoverable state

#### Documentation

 - [`b348922`](https://github.com/deis/deis/commit/b348922396958fc003b5c2df5c42e88eb8713d5c) upgrading-deis: detail deisctl versions, CoreOS security upgrade
 - [`de5e87d`](https://github.com/deis/deis/commit/de5e87d0afb29274a82ebecc865c08a255ca0d59) managing_deis: recommend wildcard CNAME or Route53 for EC2 DNS
 - [`8427d87`](https://github.com/deis/deis/commit/8427d87f57db8b4fe5cc6e1e0fc24902d2990497) managing_deis: provide login instructions
 - [`a02102d`](https://github.com/deis/deis/commit/a02102d2f4e4907dc52ff7469b8a8566f566c317) deisctl: add examples and sshPrivateKey=path to config usage
 - [`f9363f0`](https://github.com/deis/deis/commit/f9363f0d155de65b6404ec404e5ef5b8920044d8) baremetal: fix typo
 - [`d2a0046`](https://github.com/deis/deis/commit/d2a0046cf34968e6b1797c50cdefd21690215fcb) LICENSE: update copyright notices for 2014
 - [`93495aa`](https://github.com/deis/deis/commit/93495aa865f690a343cebfbfb0940e7ac157cea5) using_deis: rephrase docker image deploys
 - [`6c9a6b0`](https://github.com/deis/deis/commit/6c9a6b0952bf411ed1cc6768ba5f271a2d351bb2) client: remove mention of PYYaml
 - [`f71c6e7`](https://github.com/deis/deis/commit/f71c6e78c7d72191cc8372fbc5bc3aeb8d374e20) managing_deis: add router firewall docs
 - [`c57fd46`](https://github.com/deis/deis/commit/c57fd46dabe3267bded83d2e765b245fc05c70e4) using_deis: Process Types and the Procfile
 - [`a563d89`](https://github.com/deis/deis/commit/a563d8985ac76cccee89672bd2b3617532b1113e) manage-application: fixed bad example command
 - [`e548327`](https://github.com/deis/deis/commit/e548327277f18e684da6cd88af5aa76d61d0350b) managing_deis: always upgrade CoreOS hosts in serial
 - [`5912264`](https://github.com/deis/deis/commit/59122646c772830d0747719fb52dbdc8a88d6eb6) managing_deis: add warning to CoreOS upgrade
 - [`40d12fd`](https://github.com/deis/deis/commit/40d12fd83f3bb597a335cbd917dd1a4532825af3) contributing: clean up insecure-registry language
 - [`6c64027`](https://github.com/deis/deis/commit/6c64027ad61e7b6819da3e0983842d5ee6760c0d) installing_deis: add private network address requirement
 - [`bd14012`](https://github.com/deis/deis/commit/bd14012c180d0389c4f6a2ba5e741351e8720a12) manage-application: reference "deis perms", not "sharing" shortcut
 - [`d54ca25`](https://github.com/deis/deis/commit/d54ca25d7f3b942d4c88e0b791ad7ff1833a6df3) troubleshooting: add "could not find unit template" section
 - [`dbaf4ec`](https://github.com/deis/deis/commit/dbaf4ec7bc8dcc2bce13351fbba16cae866eafb4) troubleshooting: use double backticks for inline literals
 - [`9980e60`](https://github.com/deis/deis/commit/9980e6057f4bfb55b3bad1caf85b68f5214cf4cf) install-deisctl: remove extra installer links and simplify doc
 - [`f6aba9e`](https://github.com/deis/deis/commit/f6aba9eae489ffb192926a7164c03f68f398d4c0) DO guide: added missing ".sh"
 - [`a8bb667`](https://github.com/deis/deis/commit/a8bb6675429006d27859e5cd7b5c1c676ed8b018) schedule: add release schedule & planning documentation
 - [`56f3c7f`](https://github.com/deis/deis/commit/56f3c7f4430b986c916eb9860b5ab840a558636a) contrib/ec2: fix up stack name customization
 - [`554fdc4`](https://github.com/deis/deis/commit/554fdc415ec40f16534237894df51fb86b87bf7f) (all): clarify local cluster settings
 - [`a70f599`](https://github.com/deis/deis/commit/a70f59901f9dc910087bf1772cc9e4cc835d7945) hacking_on_dies: deisctl link fixed
 - [`442a23b`](https://github.com/deis/deis/commit/442a23ba20bb228e6cde9b68a66f13437948564b) quick-start: add "Get the Source" instructions
 - [`66fb6ad`](https://github.com/deis/deis/commit/66fb6ad9004a5ca21567763b99184dfc1df223e0) releases: add instructions for starting CI jobs

#### Maintenance

 - [`4b23b6c`](https://github.com/deis/deis/commit/4b23b6cdd829805bc3fdbca575e75ce78adcd77b) release: update version to v1.1.0
 - [`7301c59`](https://github.com/deis/deis/commit/7301c59d213a82a601440ba73f34759e9634d576) Makefile: update development registry image to v0.9.0 tag
 - [`3001ab0`](https://github.com/deis/deis/commit/3001ab05f67e32a13f113d05f3ddcc5c0b054cb9) release: bump controller API to v1.1
 - [`b0f2a6a`](https://github.com/deis/deis/commit/b0f2a6a1af8459d9301ca706e3bd41b641d858b3) README: update release version in master to v1.0.2
 - [`0fa851f`](https://github.com/deis/deis/commit/0fa851f231fdb58e935f4d4854e66c2c0eccbff8) release: update release version in master to v1.0.2
 - [`da6e2ec`](https://github.com/deis/deis/commit/da6e2ecccc48190a64a8b13598351998e7b776f1) release: bump version to v1.1.0-dev
 - [`4bcae93`](https://github.com/deis/deis/commit/4bcae93917c5407d0d155914c84952bc89ccae29) builder: bump Docker to 1.3.2
 - [`fbb7b35`](https://github.com/deis/deis/commit/fbb7b35f245d82fe3717d82c349bfde472e070ff) (all): bump CoreOS to 509.1.0 for Docker vulnerability
 - [`e1b0c2e`](https://github.com/deis/deis/commit/e1b0c2e1375743a579d66b28e35f8f9c6a462ca0) registry: Upgraded registry to 0.9
 - [`e102565`](https://github.com/deis/deis/commit/e1025656ec849acc506c376d876f537c17f3fc92) release: update version in master to v1.0.1+git
 - [`6fb1f99`](https://github.com/deis/deis/commit/6fb1f99a9882db0f30e94b0292cd6e10b27267ba) controller: update docker-py to 0.6.0
 - [`b1a3777`](https://github.com/deis/deis/commit/b1a377771b832c017b9cad52b0feea2975753e7f) (all): Added editorconfig

### v1.0.1 -> v1.0.2

#### Fixes

 - [`02dacee`](https://github.com/deis/deis/commit/02dacee1e49c88e1497e5b7c62672743414a7dc7) builder: explicitly call out $TMP_DIR on `docker cp`
 - [`3502595`](https://github.com/deis/deis/commit/3502595c3310e982f81f0914fe04b6ee0d96e0d1) contrib/coreos: fix usage of --insecure-registry
 - [`c3cede8`](https://github.com/deis/deis/commit/c3cede829ea4256c22130d08d7b3a22cfc402676) contrib: add other RFC1918 private addresses
 - [`eb70b77`](https://github.com/deis/deis/commit/eb70b772825f597bff73455ed6cb88e5a1fd580f) contrib: add --insecure-registry flag for docker 1.3.2

#### Documentation

 - [`3f1b7f4`](https://github.com/deis/deis/commit/3f1b7f475f4bbfe50c625b744c32715d39fbf8d0) managing_deis: add warning to CoreOS upgrade
 - [`67c1dfb`](https://github.com/deis/deis/commit/67c1dfb29fcb4cb1725a196f0a3a91b6a84ee599) contributing: clean up insecure-registry language
 - [`10af173`](https://github.com/deis/deis/commit/10af1737f1acc3702d04141e10093f0ba5654583) installing_deis: add private network address requirement

#### Maintenance

 - [`7940582`](https://github.com/deis/deis/commit/79405825588b03b1705de9a071242f3b60d57060) release: update version to v1.0.2
 - [`df93f71`](https://github.com/deis/deis/commit/df93f7150501dcee1ca82057de010d039fe7207a) builder: bump Docker to 1.3.2
 - [`bdc1d72`](https://github.com/deis/deis/commit/bdc1d722016bd42a0a1fae932b175296348987fa) (all): bump CoreOS to 509.1.0 for Docker vulnerability

### v1.0.0 -> v1.0.1

#### Features

 - [`90f6ef5`](https://github.com/deis/deis/commit/90f6ef518e039cf1d86884bbb7f7ec885bf8682e) security: add custom firewall in platforms without security

#### Fixes

 - [`ff9f507`](https://github.com/deis/deis/commit/ff9f5077ff3c4e22e725c67b3f78f470c0d1a31e) client: _logger doesn't exist for settings
 - [`b475c50`](https://github.com/deis/deis/commit/b475c507d15569c6330613c38e6f7f73aa915ea6) docs: update links to Docker documentation
 - [`9e2bc36`](https://github.com/deis/deis/commit/9e2bc363f65a37439b866e0dd18afd47b1920a7a) controller: change timezone data to UTC
 - [`83afc3f`](https://github.com/deis/deis/commit/83afc3f5ac20904a5ad6dc9f3d871992b9a6340d) router: Increased connect timeout

#### Documentation

 - [`efa5f31`](https://github.com/deis/deis/commit/efa5f31c9f6d7647d0cb68b783943f4d61677c24) managing_deis: add note for public cloud environments
 - [`f396344`](https://github.com/deis/deis/commit/f3963443e5ca1690d8cadd5c2e41b5c42d96eb6f) database: add BUCKET_NAME Environment Var
 - [`9b43534`](https://github.com/deis/deis/commit/9b435346e5edbd21423505176c857db2e4ffbf88) concepts: add missing word
 - [`6ad1f08`](https://github.com/deis/deis/commit/6ad1f08c84e10403d292c820cbe30e23a6723654) installing_deis: add deisctl version check
 - [`1940eaa`](https://github.com/deis/deis/commit/1940eaaafbe247b37bbae5f5a829b8f0c9870e4e) installing_docs: add cluster size to system requirements

#### Maintenance

 - [`37711ed`](https://github.com/deis/deis/commit/37711edf676bab7776de5eee40476709c643dedb) release: update version to v1.0.1
 - [`8932a0b`](https://github.com/deis/deis/commit/8932a0bd6e16177e82808625b579fb9bdfb95b7b) release: update version in master to v1.0.0+git

### v0.15.1 -> v1.0.0

#### Features

 - [`afaeb82`](https://github.com/deis/deis/commit/afaeb8218423779fd083c501c9f05d07a37f8322) (all): Do not commit a local vagrant settings
 - [`d34146a`](https://github.com/deis/deis/commit/d34146a980c3da00a29c94625fb5182041beaf43) provision-rackspace: added optional parameter to specify environment
 - [`97ce5b2`](https://github.com/deis/deis/commit/97ce5b25ca17eb1e90a5a3e4c21745b1c22c7ca5) client: optional drink of choice
 - [`270375d`](https://github.com/deis/deis/commit/270375d3d75f048f6497d0ca86c8c5c08ad6f6a5) client: Add error message when writing settings fails.
 - [`7125060`](https://github.com/deis/deis/commit/7125060e2fc1bb898c165ed6bacd57531397475f) coreos: use custom image in toolbox
 - [`47432c3`](https://github.com/deis/deis/commit/47432c3e647d4f300376a186772f508a777258ef) publisher: use godep
 - [`dd1a3e1`](https://github.com/deis/deis/commit/dd1a3e100cb04fb8b4933f27c5401638b3941fbd) logspout: use godep.

#### Fixes

 - [`42a7465`](https://github.com/deis/deis/commit/42a746557907483ad0eb7a3265fc2ef37d392a38) builder: avoid return an error if is not possible remove an image
 - [`50eb2e7`](https://github.com/deis/deis/commit/50eb2e7a868f068fdc48904242e9bd8de49bf1cb) docs: fixup markdown
 - [`62c628f`](https://github.com/deis/deis/commit/62c628f81443dc0ac2050575814b7fbe9314ee91) publisher: do not publish older app versions
 - [`b794191`](https://github.com/deis/deis/commit/b794191b21ba5e0a62536f59510a18c66c7e0f2b) docs: use /deis/logs/host|port, not /deis/logger/
 - [`dd1e3f3`](https://github.com/deis/deis/commit/dd1e3f3753564a0f2f9ba9b1bfffc8e8c6d56ee7) (all): Corrected required cores version
 - [`f26c973`](https://github.com/deis/deis/commit/f26c973779cc5999158085b3f041bbe047bdd749) router: only add x-forwarded-proto on https
 - [`2c622b5`](https://github.com/deis/deis/commit/2c622b506c7b30ed6500a7325e4a86999c4d71a8) deisctl: make restart call start and stop properly
 - [`3763748`](https://github.com/deis/deis/commit/37637487e70d7141d26186fe5a9df76e817c321e) router: remove store gateway body limit
 - [`d80ea91`](https://github.com/deis/deis/commit/d80ea91712c6afc2c9903a03e19c6d62ba22742c) client: dump log line if we fail to parse tag
 - [`c6224cc`](https://github.com/deis/deis/commit/c6224cc7ade854aac2b46ec96d41b1173baa4b50) docs: fixup links
 - [`2ed4a84`](https://github.com/deis/deis/commit/2ed4a84f5870e27c4ecc80ef2991cd49c5145f8d) builder: force newline
 - [`c0cbde8`](https://github.com/deis/deis/commit/c0cbde8bbe53df5ce04b690d60379948a8d01b60) deisctl: handle "deisctl help <command>" consistently
 - [`1da3d19`](https://github.com/deis/deis/commit/1da3d19a9c961efa13837ac5e5dde0ef196fdbad) logspout: bump packet size to 1048576
 - [`69117bd`](https://github.com/deis/deis/commit/69117bdd5fe5da5831a8f3b4878e332ab392313b) builder: remove container only if $JOB is defined
 - [`cbf38e5`](https://github.com/deis/deis/commit/cbf38e500b66716438a1954a4e7e1f528d842e04) deisctl/units: mount deis-store only if it isn't mounted
 - [`2d1fd03`](https://github.com/deis/deis/commit/2d1fd0363f60b8a463bc2837e02081974dfbd3c4) store: remove pg_num reconfiguration"
 - [`833bfd7`](https://github.com/deis/deis/commit/833bfd75a21ebd0bc4ad3b20193a62372e8c8d36) registry: prefix "apt-get install" with "apt-get update"
 - [`3056cd9`](https://github.com/deis/deis/commit/3056cd94f747178aae3f0590768f02feeb98d18a) build.sh: fail on errors in Dockerfile two-tier builds
 - [`7a4ac38`](https://github.com/deis/deis/commit/7a4ac3836f1e93fb884a0cd6f2252c3ecd21dc6a) tests: Allow TLS
 - [`dbcb2a6`](https://github.com/deis/deis/commit/dbcb2a64bae65713c727c865828516089aadd67d) logspout: remove binary from git
 - [`8c6c37b`](https://github.com/deis/deis/commit/8c6c37b509feeb0a0da73a9388cfd66fbae3e1a6) client: return lowercase app name

#### Documentation

 - [`2af6c3e`](https://github.com/deis/deis/commit/2af6c3e8449aff00a921c3e44539c951db3421e2) (all): add Deis resource requirements
 - [`fdbabb3`](https://github.com/deis/deis/commit/fdbabb33567476c2d8e583ea2a061d9f025a0ad5) README: remove "open an issue for new providers" language
 - [`3812b30`](https://github.com/deis/deis/commit/3812b30dcce1355aa7467fe7528dd31eafe6588d) readme: update readme in prep for stable
 - [`1d08e5a`](https://github.com/deis/deis/commit/1d08e5a7910f4f806d876c7d62caa38a54786a44) (all): doc edits for stable release
 - [`5e05482`](https://github.com/deis/deis/commit/5e05482dfb2dab66c7118f39ac9eb7a1c8128f47) managing_deis: add store to backup/restore docs
 - [`6efb5ef`](https://github.com/deis/deis/commit/6efb5efa8ac4854c92a9dd9c50d338afca22f609) customizing_deis: create new section for customization
 - [`cbc8bf4`](https://github.com/deis/deis/commit/cbc8bf414a97c66caa5336535a31a9c8d7282719) understanding_deis: refactor architecture, concepts, components
 - [`205f7a6`](https://github.com/deis/deis/commit/205f7a6d4ee969fce1716fd7a13bdfa3d9ea9fe9) installing_deis: add quick start guide
 - [`021e732`](https://github.com/deis/deis/commit/021e73252b3cd261cef84c4c81b29be2d5b6b497) using-deis: add missing sections on limits and tags
 - [`3ad8f9b`](https://github.com/deis/deis/commit/3ad8f9b4c69643515ba15a9f27d9607f33bc4868) test_plan: add a formal doc describing QA strategy & scope
 - [`797d3ae`](https://github.com/deis/deis/commit/797d3ae33820905df3860b42aea55341aa66d5eb) (all): flesh out platform logging docs; add logspout docs
 - [`c0228e5`](https://github.com/deis/deis/commit/c0228e5170885a22db88ec4e00d3cceeb03bb44d) managing_deis: flesh out platform monitoring docs
 - [`b87a7ae`](https://github.com/deis/deis/commit/b87a7ae9e56873dc5b712bdc57c302f2933fb481) install/upgrade: highlight that deisctl must match Deis' version
 - [`05991be`](https://github.com/deis/deis/commit/05991be548def382f3b7375e984fc62ba52113f7) installing_deis: added reference to optional environment parameter
 - [`15c4a94`](https://github.com/deis/deis/commit/15c4a94d97e502975323c15631637b6e7b6d1dd0) installing_deis: fixed formatting issue
 - [`fe4ce5b`](https://github.com/deis/deis/commit/fe4ce5bb57b8cfa07755884aa90ec676a79dd27a) upgrading: mention CoreOS minimum requirement
 - [`8434d05`](https://github.com/deis/deis/commit/8434d054d5269c44fe24c95d60c1dc951ef8759d) installing_deis: import vagrant documentation
 - [`f1264ac`](https://github.com/deis/deis/commit/f1264acaf4fdd7c9764d7d08b18796354125f8c3) troubleshooting: added instructions on manually updating CoreOS
 - [`596ec95`](https://github.com/deis/deis/commit/596ec95b3fb7bb0016d8a873be14129062a21b4a) installing_deis: add load balancer docs
 - [`3513370`](https://github.com/deis/deis/commit/35133703d02f06e94c6e437aaad35a036deddf61) installing_deis: import rackspace provisioning docs
 - [`d188ecc`](https://github.com/deis/deis/commit/d188ecc19ee44f6ff4a50cbfe79f06da4b4c07be) hacking: note that Docker 1.3.1 may need --insecure-registry
 - [`7c05e78`](https://github.com/deis/deis/commit/7c05e7817004392f8bb07802a3692e1028ce1326) testing: remove section about disabling TLS in boot2docker
 - [`9e266cb`](https://github.com/deis/deis/commit/9e266cb2bc9174c8ff4ced4e6ea181e3f52d9329) managing_deis: re-issuing auth tokens
 - [`c45348b`](https://github.com/deis/deis/commit/c45348b835ab00face8c43b04f79f7e2aca8840d) (all): specify required CoreOS version
 - [`e3d97ea`](https://github.com/deis/deis/commit/e3d97ea04ee1b19f8b57521c5d08ef47bb3e6918) managing_deis: add MDS removal instructions
 - [`32735f0`](https://github.com/deis/deis/commit/32735f01970cb1edad4c3db25c5a53c8bc9980f5) installing_deis: import gce provisioning docs

#### Maintenance

 - [`1b66674`](https://github.com/deis/deis/commit/1b66674ea606a723d7a667d50ccaa8ab5c96b97b) release: update version to v1.0.0
 - [`0fd83c1`](https://github.com/deis/deis/commit/0fd83c1714231cc745abbd71fd99f4f77753b2cc) (all): update CoreOS to 494.0.0
 - [`33b27d9`](https://github.com/deis/deis/commit/33b27d94e16da80ad2218823a623fc7681280a12) controller: update South to 1.0.1
 - [`4c63670`](https://github.com/deis/deis/commit/4c63670b3355e2c20ae89eab736003e0d7868fa8) controller: update djangorestframework to 2.4.4
 - [`08f5bf1`](https://github.com/deis/deis/commit/08f5bf1d0336cbed3bc84fc89f615dea7c7fd7a3) contrib/rackspace: bump CoreOS to 490.0.0
 - [`793a880`](https://github.com/deis/deis/commit/793a880163d0d0506614ef951bba04fd3fcfc640) provision-rackspace: bumped coreos image version
 - [`9fc6637`](https://github.com/deis/deis/commit/9fc663722420f289446dbab275592b9f68f10716) store: bump Ceph to "giant" release
 - [`38185c3`](https://github.com/deis/deis/commit/38185c3cb5c4058c2625d77777e29567e3ee5e6a) (all): bump CoreOS to 490.0.0
 - [`217261f`](https://github.com/deis/deis/commit/217261f969016033c905dccdf3e270c036f7e365) release: update version in master to v0.15.1+git

### v0.15.0 -> v0.15.1

#### Features

 - [`62d317c`](https://github.com/deis/deis/commit/62d317c03aee1cd06ac170da1e750d1e303b0956) docs: add digitalocean guide

#### Fixes

 - [`1e5b048`](https://github.com/deis/deis/commit/1e5b048f004ac310609c13838d558e0f53a8af69) contrib/userdata: nse uses `docker exec`
 - [`a5457ee`](https://github.com/deis/deis/commit/a5457ee249d8ddfd9ed48594fe1ec507c17b204f) deis/client: `Exception.message` is deprecated for >= python2.7
 - [`a139fa3`](https://github.com/deis/deis/commit/a139fa376d2a0966f63e5674569a746ea1d5a6b7) (all): etcd_set_default and etcd_safe_mkdir raise some errors
 - [`59f2bb5`](https://github.com/deis/deis/commit/59f2bb5f6c5c02219fe5ad3303c08d227f1a4152) deisctl: start gateway before volume
 - [`0ee19ee`](https://github.com/deis/deis/commit/0ee19eec0e83df2daf67f7175179230b963731c3) docs: turn off "warnings as errors" for sphinx-build"
 - [`e021b7d`](https://github.com/deis/deis/commit/e021b7df9fb76c63e5102dcf1d345be95f93a11e) deisctl/units: remove extra commands in store-volume Start-Pre
 - [`90ec427`](https://github.com/deis/deis/commit/90ec427491fb224874599edc459b3d5fae0fe8d8) docs: move subheaders down a level for correct nav
 - [`ad1c1cd`](https://github.com/deis/deis/commit/ad1c1cd0e61dac6cbdacf50c318f1a14242f7e47) docs: turn off "warnings as errors" for sphinx-build
 - [`ebf59a8`](https://github.com/deis/deis/commit/ebf59a84ba69efebd50301fe71b0e6d38ca35637) deisctl: correct typo in store-volume unit name
 - [`4dd596c`](https://github.com/deis/deis/commit/4dd596c0157eff9e31323d4232b8066424ae9aae) docs: restore "make zipfile" target for pypi docs
 - [`5aac5a0`](https://github.com/deis/deis/commit/5aac5a044392f38067fcc613841cceeb336236cf) docs: add internal links for installation guides

#### Documentation

 - [`6d43571`](https://github.com/deis/deis/commit/6d4357196ba8a8891c799b0d2d6ec7987ac7a2d1) (all): add Troubleshooting Deis docs
 - [`c0b7a32`](https://github.com/deis/deis/commit/c0b7a32c822c56afbd9a571f7853121f58bc0e8a) testing: add better instructions for contributors
 - [`5b5deac`](https://github.com/deis/deis/commit/5b5deacc1a499329d80eeae0822b8c212f1c3d0f) contrib: refer to digitalocean guide
 - [`9b36518`](https://github.com/deis/deis/commit/9b365183983fe7bd11705375ea5bd9f8857fab85) installing_deis: import baremetal docs
 - [`1e0c1a5`](https://github.com/deis/deis/commit/1e0c1a5fd93b3cb722d78459530943fbe5e84320) contrib: replace moved DNS links
 - [`4fad679`](https://github.com/deis/deis/commit/4fad679ddc61350b7229d12e3ea4a1e08f694bd5) README.md: replace moved DNS link
 - [`1864451`](https://github.com/deis/deis/commit/1864451e974770759832b476ac6501485e6349c5) deisctl: reword comment
 - [`7347ae9`](https://github.com/deis/deis/commit/7347ae9df314369476cde7bfe564d5336e282c56) README: add badge pointing to latest docs
 - [`9a614e1`](https://github.com/deis/deis/commit/9a614e1fd434e0604088bd1b6c1e70b89a9d4046) installing_deis: refer to contrib/ec2
 - [`8514a0e`](https://github.com/deis/deis/commit/8514a0ea4002d8df02a913c92dd55c565a283297) installing_deis: import AWS provisioning docs
 - [`ff3bc66`](https://github.com/deis/deis/commit/ff3bc6660bf8d81c106d00399e3674b8365ec297) reference: add controller API documentation
 - [`7517848`](https://github.com/deis/deis/commit/7517848634c2169eaead7cc64f851a5b8147ee03) (all): move installing deisctl guide to "Installing Deis"
 - [`023b630`](https://github.com/deis/deis/commit/023b6301db3bb34b3149ddf3f157f85ce97e1e0a) README: remove deprecated clusters troubleshooting
 - [`5fc3697`](https://github.com/deis/deis/commit/5fc369723c22bf6ccd00eb7134c18ae55d0fd647) install_deisctl: import shield badges
 - [`163601e`](https://github.com/deis/deis/commit/163601e1800fa1e6e83d66f39e613399a192a5a7) managing_deis: add workaround for store component failures
 - [`a8eb697`](https://github.com/deis/deis/commit/a8eb69795b442693eec340c4828314a7da370005) managing_deis: always stop platform before uninstalling

#### Maintenance

 - [`b7dd0d4`](https://github.com/deis/deis/commit/b7dd0d4fc45925de209ac01bebe3b54c220ae915) release: update version to v0.15.1
 - [`78c207f`](https://github.com/deis/deis/commit/78c207faf2625f366b5a9322b6ec024323439010) controller: update python-etcd to 0.3.2
 - [`57e650a`](https://github.com/deis/deis/commit/57e650a1091aedd4ed3dda592f38ce111b3fafc5) builder: update requests to 2.4.3
 - [`dab7a94`](https://github.com/deis/deis/commit/dab7a94461835d9b686824db72e61c08fe4dcf41) client: update requests to 2.4.3
 - [`4e8cce3`](https://github.com/deis/deis/commit/4e8cce3d9421ef782635e36214d76c8d7d49148e) release: update version in master to v0.15.0+git

### v0.14.1 -> v0.15.0

#### Features

 - [`1f74eeb`](https://github.com/deis/deis/commit/1f74eebd99d043b41d6474088cb0045712733ea0) router: add X-Forwarded-Proto
 - [`4f8bc3c`](https://github.com/deis/deis/commit/4f8bc3c2aa9be3576c9b0f0afc7f8640bb4f7df3) router: enable spdy
 - [`fd54071`](https://github.com/deis/deis/commit/fd54071f89f001b7fc47ab84e7ca5e30dd370af4) router: add optional TLS support
 - [`0cfe49a`](https://github.com/deis/deis/commit/0cfe49aea52d93a978f74ff115c61c77c963e6cb) store: add store-volume and store-metadata
 - [`2389815`](https://github.com/deis/deis/commit/2389815e849af302558d52257ba9150be780bfaf) controller: add CORS headers to api
 - [`87fe9d7`](https://github.com/deis/deis/commit/87fe9d79208445724cdd44714968ff6ff835eebc) controller: add "deis auth:passwd" to update password
 - [`241204d`](https://github.com/deis/deis/commit/241204d69d719294f3126eb6532dc3dde27b2aaa) contrib/coreos: add debug log generator

#### Fixes

 - [`6289188`](https://github.com/deis/deis/commit/6289188091f032526617a3071accce0ab9ddb1d1) controller: handle partial deletion of domains
 - [`40ead56`](https://github.com/deis/deis/commit/40ead565472dc8f61bec8a7e16e2cc4bcebab800) deisctl: "deisctl scale router=N" also starts units
 - [`d2ad06c`](https://github.com/deis/deis/commit/d2ad06ccacb2c90c274309a435d200eac0673316) deisctl: remove logspout from data plane
 - [`e4492f1`](https://github.com/deis/deis/commit/e4492f14a2f4d968eb2b61743a8d15d2fa11918a) deisctl: adjust start order for logger/store
 - [`669893a`](https://github.com/deis/deis/commit/669893ad731cf03aa21661f705f3b3ec6098f1b8) Vagrantfile: require user-data when provisioning
 - [`62ddd1c`](https://github.com/deis/deis/commit/62ddd1ca2b4621dc2e47aad7f469961504844438) tests: remove orphaned test containers in cleanup
 - [`74afcea`](https://github.com/deis/deis/commit/74afcea9d53136234aacf3cea0e24e76132903fe) deisctl: "deisctl scale router=N" uses async interface
 - [`d2d270e`](https://github.com/deis/deis/commit/d2d270e4c4eb206cec8040d3067fa6b3145da6d4) router: use string instead of JSON
 - [`a7ea4a0`](https://github.com/deis/deis/commit/a7ea4a07b4e45284f333d177865cfc322fa8658a) router: dump host/port info into json object

#### Documentation

 - [`c385735`](https://github.com/deis/deis/commit/c3857350a3b7c3d831fb46d04e3935e6103aa9e0) router: add SSL documentation
 - [`0b60803`](https://github.com/deis/deis/commit/0b608038ea9a235f586f70e41e88d07121c84ecc) releases: update with post-release "+git" procedures
 - [`44bb502`](https://github.com/deis/deis/commit/44bb5028ecb222ee32dde0739a7d80aa35c847da) README.md: remove instructions for placing apps on a different cluster
 - [`5a3bf5a`](https://github.com/deis/deis/commit/5a3bf5a54d7a018d1976e370954f04b61349e03e) router: update router's published keys

#### Maintenance

 - [`84c0692`](https://github.com/deis/deis/commit/84c0692e795f91b0e40e2975eb212f96c5250e11) release: update version to v0.15.0
 - [`9fa7c8a`](https://github.com/deis/deis/commit/9fa7c8ac81768739a0e8f3b230a1fc03aabdff76) controller: update djangorestframework to 2.4.3
 - [`354ab28`](https://github.com/deis/deis/commit/354ab285ea3d9cc23ac7cd71950ed79ca5d1e94e) controller: update Django to 1.6.8 bugfix release
 - [`1a9494d`](https://github.com/deis/deis/commit/1a9494dbd5874f28ae8c892b6bffd8ff6f6ff9d7) release: update version in master to v0.14.1+git

### v0.14.0 -> v0.14.1

#### Features

 - [`e693b00`](https://github.com/deis/deis/commit/e693b0065e2134fca9dcb9109b44e6a77dacad6a) deisctl: check for required configuration on platform install

#### Fixes

 - [`356967e`](https://github.com/deis/deis/commit/356967e859b1d0801a5c3438666c5b1806428968) tests: update smoke test for removal of clusters

#### Documentation

 - [`78518c0`](https://github.com/deis/deis/commit/78518c0106a1a31acb04dcef4662f1836b528409) (all): set required platform config before install/start

#### Maintenance

 - [`43a1dd5`](https://github.com/deis/deis/commit/43a1dd5281c0b84fc0916de74eedf0492cb48db3) release: update version to v0.14.1
 - [`a5459ed`](https://github.com/deis/deis/commit/a5459ed12174a68d69b2647687eab5249ced8832) release: update version in master to v0.14.0+git

### v0.13.1 -> v0.14.0

#### Features

 - [`85c4d07`](https://github.com/deis/deis/commit/85c4d070648f70c7d007e9ee1c1c3fc487923496) controller: bump API to v1
 - [`0b21d49`](https://github.com/deis/deis/commit/0b21d49e961ed102c520956bb5d875ebe7d37631) Vagrantfile: use "virtio" network devices by default
 - [`8881a7b`](https://github.com/deis/deis/commit/8881a7bafed270b21faeb83f5b3ed31b5652647c) client: add auth:whoami
 - [`90fb6e1`](https://github.com/deis/deis/commit/90fb6e1045a6d93ed0d9a2accfe7a1d5f3d28ec1) registry: image without development libraries to reduce the size.

#### Fixes

 - [`e61082d`](https://github.com/deis/deis/commit/e61082df0b25c48e4739183f0a6b9d215dbf0917) client: search for tag
 - [`3894f81`](https://github.com/deis/deis/commit/3894f81f344904e0ad79a9bcfa02c1603be74496) docs: update requirements in sync with controller
 - [`bca8724`](https://github.com/deis/deis/commit/bca87240702de0eacf41cea3d2af12b75a12bcc4) controller: ISO8601 datetime compliance
 - [`7bcf6f1`](https://github.com/deis/deis/commit/7bcf6f191037d797a5a0a7210873efe889600f10) logspout: add timezone to logs
 - [`48c7635`](https://github.com/deis/deis/commit/48c763531437c8ce6d5e3348859c6d33ca3ac956) client: remove timezone parsing"
 - [`87460b8`](https://github.com/deis/deis/commit/87460b869537f556d92d3ae9a7304739918a0909) controller: add timezone to datetime format
 - [`7e658c4`](https://github.com/deis/deis/commit/7e658c406dbe4709cf05ccaea70d6ca23f9f6620) builder: work around 0-byte ADD layer by ignoring tar timestamps
 - [`077b079`](https://github.com/deis/deis/commit/077b079bb606aa67c142c01dcfeeece2cc39db4c) tests: use streamOutput for test runs
 - [`77c8c18`](https://github.com/deis/deis/commit/77c8c1832a7af569d3a13f0279fe2206a7155943) client: re-parse docopt
 - [`822f1d8`](https://github.com/deis/deis/commit/822f1d83262513b4186b7d229fde530a19491343) logger: remove logger-build container after "docker cp"
 - [`0b9bd3b`](https://github.com/deis/deis/commit/0b9bd3bce8b16206ea9a8cd7cca3610ad840600d) controller: deploy only when Build is present
 - [`c69151a`](https://github.com/deis/deis/commit/c69151ab43b1d0c9711ff0e562d76abd979fea5e) builder: validate that properly formed slugs are added correctly
 - [`f2a476d`](https://github.com/deis/deis/commit/f2a476dac87fc0129a526206cf0b72b3007f3c4f) client: remove timezone parsing
 - [`0c81025`](https://github.com/deis/deis/commit/0c81025a3700aa0fa60a1c083868d0e0d5ec0ec4) controller: standardize API datetime field format
 - [`f4776f8`](https://github.com/deis/deis/commit/f4776f8ea37fbb3c40578f06efd98ec289a25f6c) builder: escape to avoid errors in serialization of Dockerfile
 - [`861b026`](https://github.com/deis/deis/commit/861b026ea28b67108bb6e524697c12aa9d08cc1c) tests: make curl loop a few times
 - [`3af7364`](https://github.com/deis/deis/commit/3af7364d207b0840028ee1e49ed0c912d2ac02a2) registry: use docker cache
 - [`4b536a7`](https://github.com/deis/deis/commit/4b536a7bd80bd8e12135f11efd521bb380222e61) controller: revert requests retry
 - [`387dfa0`](https://github.com/deis/deis/commit/387dfa0f67990614bee54ceb32e0275a4c571bd8) deisctl: optimize platform start order
 - [`8d647dc`](https://github.com/deis/deis/commit/8d647dca0658744f807c1ad9d23be32db271f439) controller: safe_mkdir /deis/domains as part of init
 - [`14e0ad9`](https://github.com/deis/deis/commit/14e0ad9c771933031b010b93a5b060aeaac8705a) (all): remove deprecated @1 unit names
 - [`dfcfba3`](https://github.com/deis/deis/commit/dfcfba3f4087699ead05394b203a535a300af87a) controller: add default retries to scheduler
 - [`cdb75ed`](https://github.com/deis/deis/commit/cdb75ed2ab72ed41e08f3e3c5e63153fad43de02) database: reset permissions before initdb
 - [`d5f922a`](https://github.com/deis/deis/commit/d5f922a266b34b9ebc3a1078542c20f19abf814b) controller: use deis-logger for requires
 - [`143f3c6`](https://github.com/deis/deis/commit/143f3c68336c451548bba3447c0c986725d8623e) deisctl: force container removal in ExecStopPort
 - [`d9367d1`](https://github.com/deis/deis/commit/d9367d105471e2eb4949967d41af75e085e6ab9a) deisctl: ignore units that dont exist on destroy
 - [`702e9a7`](https://github.com/deis/deis/commit/702e9a773eed9d51bc1e628ef76890e623711633) database: initdb if not already initialized
 - [`c3c019c`](https://github.com/deis/deis/commit/c3c019c34a9b0899b83033948eeea13262c43d22) deisctl: dont print newlines for global units
 - [`bca8eb4`](https://github.com/deis/deis/commit/bca8eb480dfd45414bb6febed200bdc575715f60) controller: use docker cache
 - [`5ac8e58`](https://github.com/deis/deis/commit/5ac8e587582fbe5928cac90fc14cafa22b9b5393) logger: kill temporary build container in "make build"
 - [`f34b355`](https://github.com/deis/deis/commit/f34b3558544953f67bf3b89c56999dde8b351642) auth: use djangorestframework login/logout views
 - [`5c8d9d3`](https://github.com/deis/deis/commit/5c8d9d3ff9796af70f6294d1cbd677eeb97ef664) builder: do not accept variables from the client

#### Documentation

 - [`6ca1293`](https://github.com/deis/deis/commit/6ca1293a97f738ac5ab85aa9c1e53ef6f26751f3) managing_deis: fix deisctl config set examples
 - [`aee7c7c`](https://github.com/deis/deis/commit/aee7c7cc60a5426d83cb7b1aaac30824a97b028b) deisctl: update with new output
 - [`486d55c`](https://github.com/deis/deis/commit/486d55c37f7858af66dc81c5a3ae77dd21a5307f) upgrade: add in-place upgrade docs, refactor migration upgrade docs

#### Maintenance

 - [`920eaf5`](https://github.com/deis/deis/commit/920eaf5db0c1c9febc772931cdbae5ad6d27c3bf) release: update version to v0.14.0
 - [`bb5bc9e`](https://github.com/deis/deis/commit/bb5bc9e1f778748745ab2657c37d2820033a68dd) (all): bump CoreOS to 472.0.0
 - [`dadfd1b`](https://github.com/deis/deis/commit/dadfd1b9f4dfa405291af52903a6ac40287123eb) deisctl: switch data containers to ubuntu-debootstrap:14.04
 - [`a564287`](https://github.com/deis/deis/commit/a564287d071b902f416e54755ebab4a9fda58f6c) deisctl: bump godeps for coreos/fleet
 - [`d40b1b3`](https://github.com/deis/deis/commit/d40b1b31fb5496ec6ece26126dfd9116bd9f01c7) store: add start delay on OSD recovery
 - [`6ea8d57`](https://github.com/deis/deis/commit/6ea8d5759b296c9d8454baedf73c42b09eff61a6) builder: update Docker engine to 1.3.0
 - [`7ba9b17`](https://github.com/deis/deis/commit/7ba9b17a4880270660c7c9a3e0840c4e4be39aed) controller: update json-field to 0.5.7
 - [`0394cf3`](https://github.com/deis/deis/commit/0394cf3eac3f410e8f2cc4f6a279388d16090090) controller: update psycopg to 2.5.4
 - [`903721b`](https://github.com/deis/deis/commit/903721b5592054d0a4e456c6022fbbb98eeb5ed0) controller: remove unused django-allauth app
 - [`d67179b`](https://github.com/deis/deis/commit/d67179b4fe716d540565bdc838416ca1b766b177) (all): bump CoreOS to 471.1.0

### v0.13.0 -> v0.13.1

#### Features

 - [`5688c6b`](https://github.com/deis/deis/commit/5688c6b3def8a45ab76e366c30e0793660df5f82) builder: inject GIT_SHA into builder apps
 - [`2c50d5c`](https://github.com/deis/deis/commit/2c50d5c2c7bf34f26b2f6e605673a78cc3602602) client: add loading info msg to run command
 - [`91164f6`](https://github.com/deis/deis/commit/91164f6d5c6fd9ec632d505a480d84cfe3f4eb97) registry: image without development libraries to reduce the size.
 - [`da80165`](https://github.com/deis/deis/commit/da80165693661ccb0df4a3944bad024cd0043d6f) contrib/util: add script to generate project contributors

#### Fixes

 - [`91cdd70`](https://github.com/deis/deis/commit/91cdd704cb629b6686700e2a427a4b9ad4533cfe) builder: use proper dockerfile syntax
 - [`bb352b4`](https://github.com/deis/deis/commit/bb352b48449cb08831c9f9126adc8d0b77a2279d) controller: use build tag if present
 - [`8418141`](https://github.com/deis/deis/commit/8418141d64dd5eb4288f3bfecb7b7bbbe5527036) controller: do not commit to latest
 - [`b65e499`](https://github.com/deis/deis/commit/b65e4995129c7c62a6224c81a68e2823f1bac370) controller: inject config on top of existing environment
 - [`f98fea6`](https://github.com/deis/deis/commit/f98fea634784f8e8077a4be130e10302daaa3ce9) controller: don't clone run containers on deploy
 - [`2b8aa07`](https://github.com/deis/deis/commit/2b8aa07ee34d2e46730d7c0436531b355afe874c) controller: set default config owner to request.user
 - [`70da474`](https://github.com/deis/deis/commit/70da474dd73dc749f9913f87d7f1a8cb38d5556a) deisctl: control plane components should restart on failure
 - [`2084164`](https://github.com/deis/deis/commit/20841647809d178e4d1428a4b2b5ccf1c722bdbb) controller: include message to avoid confusion building the component
 - [`7f6c202`](https://github.com/deis/deis/commit/7f6c20295d4b16ad9705bfeae2af1106f337e534) builder: properly parse config vars
 - [`d63aa6d`](https://github.com/deis/deis/commit/d63aa6de73dde746a4f55b89a273ab03a1d7cd5f) controller: retrieve logs via GET request
 - [`e310501`](https://github.com/deis/deis/commit/e310501ef0ea6426e2d5ac5078e150d90ea7302b) controller: properly serialize JSONField objects
 - [`dac9161`](https://github.com/deis/deis/commit/dac91613ecb0bc7476650c91a5dad9386c4613f2) logger: the test is expecting the message "deis-logger running".
 - [`d925056`](https://github.com/deis/deis/commit/d925056a83b573bc87c03b8467911aab0d605d36) controller/tests: add fleet socket to controller tests
 - [`9bf5e8e`](https://github.com/deis/deis/commit/9bf5e8e5a5b3ced96e625b2a516ac212a1c9d5ff) registry: use YAML's nil type
 - [`250e7cc`](https://github.com/deis/deis/commit/250e7cc9e0f0aa3513cdf1dae9002172a968b748) registry: make deis-cache mandatory
 - [`1fd537b`](https://github.com/deis/deis/commit/1fd537b86324b9ad8e38d02f35dbff7ca8129bdc) README: update current version in badge
 - [`b8c40a6`](https://github.com/deis/deis/commit/b8c40a64d7c247c693c8a1e5dbf0d27e4c52001b) deisctl: update `deisctl --version` string

#### Documentation

 - [`c2eb07d`](https://github.com/deis/deis/commit/c2eb07dd5936806f3a19b4267887419a01fd286a) README: add build status badge for CI test-master job
 - [`1e7ce1d`](https://github.com/deis/deis/commit/1e7ce1df0451f0a4976ee06f9a04cfa3616fff0a) managing_deis: update controller_settings.

#### Maintenance

 - [`b9db46a`](https://github.com/deis/deis/commit/b9db46a53511f9a7f6f44d7bb5d5d44673ee528d) publisher: Add missing deploy target
 - [`c1c66f9`](https://github.com/deis/deis/commit/c1c66f97afc3b99d6f20e07619dfa410d3cc81ac) deisctl: remove redundant CHANGELOG.md

### v0.12.0 -> v0.13.0

#### Features

 - [`664b579`](https://github.com/deis/deis/commit/664b579e272bc2f557ea679c62edb78c78451d9b) router: Changed store hostname to deis-store
 - [`74a018a`](https://github.com/deis/deis/commit/74a018a214d6b4330a18d7bcf862d069ef4903e3) client: "make installer" creates distributable CLI package
 - [`e5bf847`](https://github.com/deis/deis/commit/e5bf847e651bd9cc57ae091d56815db3d59d0199) deisctl: add store support; drop database-data and registry-data
 - [`4eb0089`](https://github.com/deis/deis/commit/4eb00899ea5e58a1cab2a72f0fe54e2fc5d9870e) registry: use store component for filesystem layers
 - [`b7c2990`](https://github.com/deis/deis/commit/b7c299053abb173e2ad2a6e2af39b0169a1f9a4e) database: use deis-store and WAL-e to ship WAL logs
 - [`b5cb742`](https://github.com/deis/deis/commit/b5cb7427b3e83381afa4aa6e52c29231c5a18862) store: add deis-store component
 - [`481fe92`](https://github.com/deis/deis/commit/481fe92911253069cda1f43b91c4d9a73017e1ed) router: Improved router defaults
 - [`1280daf`](https://github.com/deis/deis/commit/1280daf770263eff1c98caca4d10efa4d61304d2) controller: introduce deis/logspout
 - [`0dcb368`](https://github.com/deis/deis/commit/0dcb3681f91536d9f12a8ecb834bb655362be1d8) controller: fix flake8 error.
 - [`0254f28`](https://github.com/deis/deis/commit/0254f28a760558eeb544587fd9607a0719661d20) controller: restart the app if there is a failure.
 - [`7df6688`](https://github.com/deis/deis/commit/7df6688d6ad92be4bbc456678d3275ef61cbcbb4) contrib: add bumpver tool to help with semantic version releases
 - [`d1066c8`](https://github.com/deis/deis/commit/d1066c8ff1eaeae02d54c2349b226152cceb1a92) tests: Show all etcd keys on error
 - [`c27bd53`](https://github.com/deis/deis/commit/c27bd53bae98abca3fddea2bee52136ee327e526) router: image without development libraries to reduce the size.
 - [`0b0cb6f`](https://github.com/deis/deis/commit/0b0cb6fc6f85fc026697cc78a3ae6f8373eb4245) deisctl: add deis/publisher
 - [`a9c98c1`](https://github.com/deis/deis/commit/a9c98c1b06a4976fdb84c90b265b0afefc919c45) cmd: refresh-units accepts -t for tag, branch, or SHA
 - [`3158d8b`](https://github.com/deis/deis/commit/3158d8bbc99f5995c052547a120993158c2be80d) userdata: add deisctl to the coreos install
 - [`ee5f8a3`](https://github.com/deis/deis/commit/ee5f8a3e8a5af58650971df1c2845eee3a800752) contrib/digitalocean: Use native CoreOS integration
 - [`938f4eb`](https://github.com/deis/deis/commit/938f4eb3c2c490d0f92dce25d1fd4f684012dfee) cmd: add informative messages to install
 - [`c45fade`](https://github.com/deis/deis/commit/c45fadeebe09026599d90193ffd594016cd7e7e7) deisctl: colorize deisctl
 - [`615a2cb`](https://github.com/deis/deis/commit/615a2cb40aa7585428f596114b4270027e1e8ab5) start/stop: allow starting or stopping > 1 unit at a time
 - [`077429c`](https://github.com/deis/deis/commit/077429cb472553ef2818c9fc27653d574203f3e7) start: wait on containers to start
 - [`74e9337`](https://github.com/deis/deis/commit/74e9337f174a928d68d33f657afa4655bce2582b) installer: use /usr/local/bin and one-liner install scripts
 - [`2ae953e`](https://github.com/deis/deis/commit/2ae953ece143a84fe891f10d039902acef28945b) Makefile: create shell script installer
 - [`6d541aa`](https://github.com/deis/deis/commit/6d541aae5dd141ea8ee3bd1434c6d233770fb37d) config: first pass at config subcommand
 - [`c182966`](https://github.com/deis/deis/commit/c1829662a1f02eadf7902963f62207f22ffaf5f4) restart: add restart command for convenience
 - [`bf68a8e`](https://github.com/deis/deis/commit/bf68a8e64d9095ddbc151ca61fae5a868646aff7) journal: add journal support
 - [`65ee91e`](https://github.com/deis/deis/commit/65ee91eaafdebf672c78d3205e7354806ecb8922) deisctl: move server env variable to etcd
 - [`975ac6b`](https://github.com/deis/deis/commit/975ac6bb840bd492d1863564691598ed91233868) deisctl:  add new feautes and update core-os updatectl to updateservicectl
 - [`b366679`](https://github.com/deis/deis/commit/b366679d4111b6b098f791ae9f559ed999c819e3) deisctl: get groupid and app id from etcd defaults to env variables
 - [`8509f01`](https://github.com/deis/deis/commit/8509f0147dcc8921b781ec438876eadbfa1b8783) deisctl: removed unneccessary code
 - [`558bfb0`](https://github.com/deis/deis/commit/558bfb0660ce1eea91625dfa8ae51f05bd928afe) deisctl:  add hooks and units dirs to constants
 - [`cdd0150`](https://github.com/deis/deis/commit/cdd015093622f30e50f34aa3837892c0ec3433b1) hook: add pre/post update hooks
 - [`74f5229`](https://github.com/deis/deis/commit/74f5229658ad61a7583e166d04cca63bccf75078) deisctl:  working version of updater latest
 - [`d96d9bb`](https://github.com/deis/deis/commit/d96d9bbc1b97bd1c58c1bd0d9b9855ccf9453ee7) deisctl:  basic working version of updater
 - [`4c509b1`](https://github.com/deis/deis/commit/4c509b1055e629aab3eb743688b8da90d5a32fb0) deisctl:  change flag to string
 - [`5a5f859`](https://github.com/deis/deis/commit/5a5f85988b9b28a9b123b6bee611d36cf6db66b2) deisctl: fix package names
 - [`d07a923`](https://github.com/deis/deis/commit/d07a923d167535380441f3c1b7e86d2502cdfcf1) deisctl: updated command instance
 - [`cb18218`](https://github.com/deis/deis/commit/cb18218d3661ae94854166df759e59853273cd1b) deisctl: removed systemd and distibuted lock
 - [`82ccd66`](https://github.com/deis/deis/commit/82ccd6689aa9180e5d075840d64ba2068efb4803) deisctl: add utils for client instance functions
 - [`dcd22f1`](https://github.com/deis/deis/commit/dcd22f1b14571b9b65bfb4a7a9e11ec0b31bd19f) deisctl:  add update command and updatectl package

#### Fixes

 - [`4049d70`](https://github.com/deis/deis/commit/4049d70dd50f157650ab4975794c119d3335625b) deisctl: escape $HOME in post-install script
 - [`99ccce4`](https://github.com/deis/deis/commit/99ccce4afa5ae7d8fc7338d4c472ce2d7cfdbae7) controller: use build image for publishing releases
 - [`be5ad23`](https://github.com/deis/deis/commit/be5ad2327cc28c481c90fd623fc95fed60febca2) deisctl/units: don't fail units if ExecStopPost fails
 - [`bc7853a`](https://github.com/deis/deis/commit/bc7853af475b6adab0f6be9b3cdd1708c90f15fc) deisctl/units: short-circuit if we don't have an image
 - [`13401a4`](https://github.com/deis/deis/commit/13401a4fa6176ecd0b8dc90d1054fd6c43b97ea4) controller: increment run containers and keep db models
 - [`02213f9`](https://github.com/deis/deis/commit/02213f979458f5c6266b0425ee3d246f9b3403e6) deisctl: bump default etcd timeout to 10s
 - [`f81c766`](https://github.com/deis/deis/commit/f81c766c652c9c5c5e0083ca0044067a5bbe0945) controller: add retry on destroy
 - [`3eb39c0`](https://github.com/deis/deis/commit/3eb39c067dd974905935c908f758608c5dbc7116) controller: work around fleet state reporting bug where units report as "failed" before they go "active"
 - [`2c1ae82`](https://github.com/deis/deis/commit/2c1ae8206af983fbf1228a122ead8dcd2754a001) controller: add retries on container creation
 - [`9e86f87`](https://github.com/deis/deis/commit/9e86f877dd15737e65506fe680b882b85111517a) deisctl: bump status pollwait sleeps to 3s
 - [`ea51772`](https://github.com/deis/deis/commit/ea517729e7e1fa08476e13a088407ba8d85ec157) deisctl: don't set SSH tunnel for refresh-units command
 - [`f6f11b0`](https://github.com/deis/deis/commit/f6f11b0c2fa14719764e1a95e68f742aa90c9437) tests: correct typo in cli installer versioning
 - [`b20d2c8`](https://github.com/deis/deis/commit/b20d2c8941f491e27d5a19bfa1a4ad1bb8aedeba) publisher: skip when failed to retrieve container
 - [`39a5e7f`](https://github.com/deis/deis/commit/39a5e7faf433a6b324fa09a6866f6558c585ee12) store: remove pg_num reconfiguration
 - [`18cad13`](https://github.com/deis/deis/commit/18cad13b1415feb72377855f10d0d5c2b423eae3) (all): resolve store rebase issues
 - [`0ceb125`](https://github.com/deis/deis/commit/0ceb1255ba32932ea275f460e1bb5529df5a5505) tests: dont use git remote on deis pull test
 - [`8a3fa67`](https://github.com/deis/deis/commit/8a3fa672198804dfa3516ec1e2bf7764268ccaf0) builder: remove vfs on tests
 - [`83e585e`](https://github.com/deis/deis/commit/83e585ea3a4be2e0109949449b1476c3a141fa20) (all): add store dependencies to tests
 - [`4eb0089`](https://github.com/deis/deis/commit/4eb00899ea5e58a1cab2a72f0fe54e2fc5d9870e) registry: use store component for filesystem layers
 - [`b7c2990`](https://github.com/deis/deis/commit/b7c299053abb173e2ad2a6e2af39b0169a1f9a4e) database: use deis-store and WAL-e to ship WAL logs
 - [`b5cb742`](https://github.com/deis/deis/commit/b5cb7427b3e83381afa4aa6e52c29231c5a18862) store: add deis-store component
 - [`29f2b40`](https://github.com/deis/deis/commit/29f2b406251fad0066a6f4d7bb81e18048007c5b) controller: alter logic on container creation
 - [`53e93e8`](https://github.com/deis/deis/commit/53e93e8535c8d1e9753787b24405a5c7a4798b08) contrib/ec2: fix HVM AMIs for CoreOS 459.0.0
 - [`f6093db`](https://github.com/deis/deis/commit/f6093db6a29578079e0773572d77cecc20cc1ff5) deisctl: refresh-units uses $HOME/.deis/units by default
 - [`2ce2fca`](https://github.com/deis/deis/commit/2ce2fca501b4243d7948ab744abe7e8163a1701b) tests: Increased test timeout for smoke tests
 - [`a2f1475`](https://github.com/deis/deis/commit/a2f1475f6a215a88b47d0607c23b8cedce337629) builder: compiled slug is no longer a tarball
 - [`dcf37ea`](https://github.com/deis/deis/commit/dcf37ea00c138fc7b72a55f010eef4b87fd39ebc) router: Use official javascript content type
 - [`5b474ae`](https://github.com/deis/deis/commit/5b474aee9d8e8dd83a9346b05e629df8fb001e63) deisctl: fix global unit display
 - [`88a6eea`](https://github.com/deis/deis/commit/88a6eea678d833b7317c7d509a6ab2f5cefaffd7) builder: do not add build option if empty
 - [`b989ab7`](https://github.com/deis/deis/commit/b989ab75d5407e189b5c83c40b2d7c122dceccf1) tests: mock postgresql should wait until initialization complete
 - [`9dcdfd6`](https://github.com/deis/deis/commit/9dcdfd64d7254dea513173b0a75c92e9affcfd92) builder: pipe envvars to temp file
 - [`59f3b7d`](https://github.com/deis/deis/commit/59f3b7d0d191117c55894c356d0da5a025a2eeaf) router: update to confd 0.5.x branch to fix sort issue
 - [`6512cfd`](https://github.com/deis/deis/commit/6512cfd9f97f464f90569df96b0847d4a9f6c672) logspout: use custom log format
 - [`e39edb6`](https://github.com/deis/deis/commit/e39edb65d74d88db5dab4c5744cfa2b1cc6bdb60) logger: trim newline character from message
 - [`4a56691`](https://github.com/deis/deis/commit/4a5669111ed0915e360d09c3d97792b82f1edc79) deisctl: add timeoutstartsec to logspout
 - [`2cdf3ee`](https://github.com/deis/deis/commit/2cdf3ee77046bda80baf33970b9f34d57102b929) deisctl: remove PostExecStop for global units
 - [`15b0843`](https://github.com/deis/deis/commit/15b0843519a04fd76f37795aa75837153af9e21d) logspout: add etcd client request timeout
 - [`a315026`](https://github.com/deis/deis/commit/a315026ab609e7ba9d334325932086e776184309) builder: suppress stderr
 - [`eb27432`](https://github.com/deis/deis/commit/eb27432391b5e6d6e39813399901c889ce920963) builder: exit if response is bad
 - [`bd384bd`](https://github.com/deis/deis/commit/bd384bda890ff709768b52615bfc019ed040fc02) tests: add registry log output on error
 - [`ad2981d`](https://github.com/deis/deis/commit/ad2981d6cd82d352e7e00660a17e4661c9bc3ede) tests: shrink deis pull image by 200MB
 - [`0b6635b`](https://github.com/deis/deis/commit/0b6635b38a37f9dc84c32a2c258b6e3ed7d52629) router: change nginx directory to /opt/nginx.
 - [`ccffb4f`](https://github.com/deis/deis/commit/ccffb4f683f76e5cc4272bd5a077513f9a8b4fd0) router: use docker capability to decompress files.
 - [`f7f00c1`](https://github.com/deis/deis/commit/f7f00c1dd9c7cc217f1aa30672b5e862ab381322) router: Set default values in /deis/router
 - [`8ddb7dc`](https://github.com/deis/deis/commit/8ddb7dc122304b09d73b4445414b61e2c80ff594) controller: Use correct entrypoint for run
 - [`ebeab96`](https://github.com/deis/deis/commit/ebeab96b2e8994488620f1757bd3ca82a2288efe) scheduler: allow proctypes with dashes and underscores
 - [`78eb4ec`](https://github.com/deis/deis/commit/78eb4ec445de9927d94a0bc5e4062c0bc36b9c52) logger: log the error, but don't panic
 - [`89c1f41`](https://github.com/deis/deis/commit/89c1f418a638d84c6dedb53759aa4b10b1f54163) logger: add more verbose message
 - [`da17824`](https://github.com/deis/deis/commit/da17824ffb790e68bf9da1de3c29f67e1c705fce) logger: don't print a message on every incoming message
 - [`aa98983`](https://github.com/deis/deis/commit/aa989830a7511d108fb91304a7ed178136afc904) builder: restore VOLUME instruction so builder does not stack overlay filesystems
 - [`b38763a`](https://github.com/deis/deis/commit/b38763aaff45de0ea789686daf8df623490ce859) router: use ubuntu:14.04.
 - [`0d114a6`](https://github.com/deis/deis/commit/0d114a65469e60c20e12ada7670b2d5f767ed92d) deisctl: typos in publisher unit file
 - [`c2119d4`](https://github.com/deis/deis/commit/c2119d4f43e1970842e92f458c8f58e1c20d92f0) router: ubuntu-debotstrap doesn't have netstat and confd errors are redirected to /dev/null.
 - [`f521001`](https://github.com/deis/deis/commit/f52100171cff438ff8d34c49333afa9e86b23ca8) router: build cleanup.
 - [`fef3f1a`](https://github.com/deis/deis/commit/fef3f1aaf4c93ad623333680049b02227abb6073) builder: Reverted to default php buildpack
 - [`edca91b`](https://github.com/deis/deis/commit/edca91bc54be2eac592629a35031afbb2095b718) builder: Non-root slugrunner and slugbuilder
 - [`c092fed`](https://github.com/deis/deis/commit/c092fed075e71e87dcce0a1103fe7cba4d65c1ea) builder: exec the runner
 - [`9c203ca`](https://github.com/deis/deis/commit/9c203ca084e579a232ece13b64e4cb76c30e8f2f) publisher: fixup entrypoint
 - [`74ef30c`](https://github.com/deis/deis/commit/74ef30c5c4e8e61281d446ec0a2607c0dbc76128) publisher: use git sha tag for `make build`
 - [`1655fa5`](https://github.com/deis/deis/commit/1655fa504206e192ea5cc640483a646baf9f5873) controller: bump timeout to 20 min for image to download
 - [`af4a0d0`](https://github.com/deis/deis/commit/af4a0d0f5ce79142a3978a5a7abf6e620a3ae69e) controller: work around fleet state reporting on deis run
 - [`f9c6828`](https://github.com/deis/deis/commit/f9c682856b0bb685ac1aa090ae3ab7daade508cb) deisctl/units: fix ncat loop in builder and registry
 - [`640892a`](https://github.com/deis/deis/commit/640892a1d9a21ebae815040add2ecb744493afc8) tests: data containers should not be shared
 - [`f8dd119`](https://github.com/deis/deis/commit/f8dd1198014a13d4dab67d3214c90819b9d926fd) makefile: use portable sed for make discovery-url
 - [`c34be6b`](https://github.com/deis/deis/commit/c34be6b63868aabe11fc13bfba0159355e9c193a) deisctl: create unit download URL properly after move to subdir
 - [`d45b4d9`](https://github.com/deis/deis/commit/d45b4d9c16c83b516ed3d2ead12f3bab07370c0f) tests: add correct key to ssh-agent
 - [`c35fda4`](https://github.com/deis/deis/commit/c35fda48f65aac5a69c6061173e8baf613a3d32c) docs: clean up sources and treat warnings as errors
 - [`e271725`](https://github.com/deis/deis/commit/e2717250843af3ba702717ac4a67279d22d4a692) tests: add DEISCTL_UNITS to test-setup
 - [`3710994`](https://github.com/deis/deis/commit/37109941286fde45d962d93d0790e3af5c915b70) makefile: replace user-data inline
 - [`4d853a7`](https://github.com/deis/deis/commit/4d853a773b447e1962162c7c44dbac29515d0931) tests: install required go dependencies
 - [`d45dc77`](https://github.com/deis/deis/commit/d45dc77dba02ba5d4df186c27ac653292c3dacf8) tests: add log_phase helper
 - [`c0887df`](https://github.com/deis/deis/commit/c0887df7417322f45c27f0b2e94f5175a6c58268) deisctl: install only deisctl
 - [`2f17744`](https://github.com/deis/deis/commit/2f17744b3485e2d5fa28e1af23a533d5bbeaefba) readme: use relative link to deisctl in repo
 - [`897bec2`](https://github.com/deis/deis/commit/897bec27c9810db90209ff5a4a6afeca53340697) CHANGELOG: use SSL for github and recreate change entries
 - [`b91abd7`](https://github.com/deis/deis/commit/b91abd7ace1580f011da836c03e9d0e3f3e6aa74) router: gzip content types syntax fixed
 - [`dfca1d1`](https://github.com/deis/deis/commit/dfca1d13bdcdf7567668d74e328b6c8aad90250e) (all): Announce the published port to etcd
 - [`e8aa998`](https://github.com/deis/deis/commit/e8aa998e1a43406b67eb408de80fdbe7d8921b11) (all): Announce the published port to etcd
 - [`85b7922`](https://github.com/deis/deis/commit/85b79224663428ac463e1e4ecc8b39f95a35a29d) deisctl: Use $HOME/.deis/units, don't expand tilde (~
 - [`322bc2f`](https://github.com/deis/deis/commit/322bc2f37d711100060d175911a843b5854042a2) Makefile: static complication in installer
 - [`742da09`](https://github.com/deis/deis/commit/742da09b30dbf9f82f9027b11103c378fee842c6) help: `deis help` lists all available commands
 - [`5933146`](https://github.com/deis/deis/commit/59331461d559552fd1daccefc6dee8cbfa80a1cf) Makefile: static binary
 - [`683d054`](https://github.com/deis/deis/commit/683d05429b146bb8b7d5e715cc12508fd5a22288) Dockerfile: deisctl is a static binary
 - [`be4bef3`](https://github.com/deis/deis/commit/be4bef39543eb5bc391899e18ac95b65627a46f7) state: ignore intermittent timeouts when polling for UnitState
 - [`e316572`](https://github.com/deis/deis/commit/e31657238b1e5a5f5b6d82e821f1300936370aa0) Makefile: ensure the installer makes /var/lib/deis/units readable
 - [`d831066`](https://github.com/deis/deis/commit/d8310669d884a7a1b7cb152415c74001229f593b) registry: fix duplicate tag in registry-data unit file
 - [`561fb31`](https://github.com/deis/deis/commit/561fb31baa434881886cb810074ddecf43c62153) units: show start-pre status when downloading data container base
 - [`81371d5`](https://github.com/deis/deis/commit/81371d59c20325051e7287326a6b9648e53d2a61) cmd: only allow the router to scale past 1
 - [`05f027c`](https://github.com/deis/deis/commit/05f027c6fc89ebc731da91ead64628c0e7e093da) cmd: create unit files as readable to all users
 - [`043739e`](https://github.com/deis/deis/commit/043739ea7f87a70cf087cf17d98ecb9c02367e09) client: destroy all units if none specified
 - [`40b084f`](https://github.com/deis/deis/commit/40b084f4384b83a79061dfc2a920649b2a0581fa) client: dramatically simplify scaling logic
 - [`37c348e`](https://github.com/deis/deis/commit/37c348ef2f94ef62b51b09381c41f9abccea85e6) cmd: append "@1" if none supplied to install
 - [`8f14e8a`](https://github.com/deis/deis/commit/8f14e8aea22f415b5e14624a4371b15b5baed327) client: destroy all units if none specified
 - [`725ae0d`](https://github.com/deis/deis/commit/725ae0d84e51e466f31cd864521a856773015011) client: check if unit exists
 - [`7bd1709`](https://github.com/deis/deis/commit/7bd17097de3e1a6d0f98d86daa95d863323d69ef) client: return error if unit list is empty
 - [`1f79c8b`](https://github.com/deis/deis/commit/1f79c8bbc6f09861b2af05a20c31745430d0e92b) deisctl: use docopt's native version parser
 - [`fd18f5b`](https://github.com/deis/deis/commit/fd18f5b77491a5950ddc56966f33a2f8eb71c308) makefile: expand paths for golint
 - [`cc9e0a4`](https://github.com/deis/deis/commit/cc9e0a49785edf0b11be97d3b75df44464415570) units: add fleet.sock bind-mount for controller
 - [`4497af2`](https://github.com/deis/deis/commit/4497af28bb2de71205623a25462fd374b2eb6180) cmd: allow `-p` to specify where to save local unit files
 - [`318514e`](https://github.com/deis/deis/commit/318514e3bbf8f11577e70b65a52a180dbcab7691) client: look for unit files in ~/.deisctl before /var/lib/deis/units
 - [`c325ca2`](https://github.com/deis/deis/commit/c325ca28c085acc5ffbec8432cec31d8e820da8e) debug: remove other vestige of unused --debug flag
 - [`72cb9d8`](https://github.com/deis/deis/commit/72cb9d89fddded0057fc6de416170236eef45163) README: fix installer link to use http, not https
 - [`54dc7df`](https://github.com/deis/deis/commit/54dc7dfc9aaea32a1f2d16897d05c90f2e56171c) README: update installer link
 - [`ba85174`](https://github.com/deis/deis/commit/ba851748bde851f61462a0010d7c5f0d7055e155) installer: use deisctl-hack fork of makeself
 - [`dce122d`](https://github.com/deis/deis/commit/dce122d2ec7a770b6f3dcb168ee591a83f6e6bc6) debug: remove unused --debug option
 - [`9459a83`](https://github.com/deis/deis/commit/9459a835089a6fa7c3fd37d5e3d8b4da22214033) version: add special handling for --version
 - [`3aaf764`](https://github.com/deis/deis/commit/3aaf764cceadb29d88499d716f669e38ef03359c) cmd: add explicit platform target
 - [`4b9f157`](https://github.com/deis/deis/commit/4b9f157dbd63b0b3589b7f4bd3d6174b851aab8a) tests: explicitly set tunnel to null
 - [`dce34ce`](https://github.com/deis/deis/commit/dce34ce3edc89f6eb12a334a2963625c46e5ab71) destroy: fix shadowing bug in destroy
 - [`3227656`](https://github.com/deis/deis/commit/3227656ad05a171b30cb63f1c8b9cdd62f619744) units: controller waits for logger container in ExecStartPre
 - [`471e4e8`](https://github.com/deis/deis/commit/471e4e87c899288ea130962c97be68103db18012) units: use @ in wildcard for router conflict
 - [`c2b75ee`](https://github.com/deis/deis/commit/c2b75ee318ed1f6adaeadd941908d64b7d89f881) unit: match @ units properly
 - [`f6e0b86`](https://github.com/deis/deis/commit/f6e0b86d21a077c8376c5eb17eca4e3a055c004d) destroy: wait for job state inactive on destroy
 - [`e6d260a`](https://github.com/deis/deis/commit/e6d260ad150d185c53928b7d20486fb49f4db2f0) ssh: switch to default known hosts
 - [`cba85ee`](https://github.com/deis/deis/commit/cba85eeff5797fb9957a22695b5e37c2792801f3) units: use default GOPATH for unit lookup, if available
 - [`a6dc5a9`](https://github.com/deis/deis/commit/a6dc5a9818a43ce3895c3d4d00c1f0484494cc12) update: fix imports
 - [`700ad7b`](https://github.com/deis/deis/commit/700ad7b1185c39369a4fa0972fe07553df89b23b) state: print inactive states without substates
 - [`265cc3d`](https://github.com/deis/deis/commit/265cc3d8c765b42c7eb23dda0e4dc4291ee1df02) units: switch to systemd template units
 - [`985f003`](https://github.com/deis/deis/commit/985f0039729bfdf8843344d4040b14668dd9b882) deisctl: fix utils error
 - [`7358f4d`](https://github.com/deis/deis/commit/7358f4d7ea7dc28d35427fd7b14769745d9314cf) update: extract update to root
 - [`1d629e5`](https://github.com/deis/deis/commit/1d629e566858456dfe6a1b34606ed38ad3358f12) update: add update service as systemd unit
 - [`d6ccce2`](https://github.com/deis/deis/commit/d6ccce20c906d1ef2ec07a4c8b8a1149910544bd) updatectl: fix data container matching, fallback to envvar for version
 - [`88745a8`](https://github.com/deis/deis/commit/88745a825d66d11caef6d588ae657268bb6f21f1) update: do not pull images on update
 - [`329c372`](https://github.com/deis/deis/commit/329c372cd153eb481da25de20803262e56fda32b) constant: add new constant package
 - [`7c2c5db`](https://github.com/deis/deis/commit/7c2c5db2c7d42c015e39d1d77ad2b372274a2518) (all): rename constant folder, go fmt
 - [`63e7db9`](https://github.com/deis/deis/commit/63e7db9f804491aebeb321b5c751634c27f099e4) units: cleanup post-start output for builder/registry
 - [`910ccef`](https://github.com/deis/deis/commit/910ccef6c735a02ff0e773f488969bf21c26aa7d) install: install registry after cache
 - [`7687c1d`](https://github.com/deis/deis/commit/7687c1d7e26fe3a2a52041d53eeba7110d1da2e1) units: switch to new fleet X-ConditionMachineID
 - [`4be9b61`](https://github.com/deis/deis/commit/4be9b6192b573f55d1fe26515a72911efaac1716) packaging: add version to package tarball
 - [`557cefd`](https://github.com/deis/deis/commit/557cefde26489774e19d483b8d76554f611cb5cf) packaging: update Dockerfile and paths
 - [`f7363fa`](https://github.com/deis/deis/commit/f7363fae83f1ee9c46c94198cc73c7d1939666bc) upstream: rebase against fleet upstream changes
 - [`5602a17`](https://github.com/deis/deis/commit/5602a17f2573c0d4ae3b4dce78bf9dbf99185006) main: fix package path

#### Documentation

 - [`11b75e7`](https://github.com/deis/deis/commit/11b75e71e5b30a7737bfabde31f663ff0d719d4b) store: remove not-implemented store-data
 - [`fe25f72`](https://github.com/deis/deis/commit/fe25f721244cb54da507e635cea0fa44419b1359) ec2: move "scale routers" section after platform installation
 - [`ed04d75`](https://github.com/deis/deis/commit/ed04d755dee51b8ad89e8318b8e46df9955068a4) managing_deis: remove dead machine from etcd cluster
 - [`46e0cfa`](https://github.com/deis/deis/commit/46e0cfad2631d10a575d9ac9b035f7ac6b121bb2) managing_deis: link to Ceph troubleshooting guide
 - [`9080b77`](https://github.com/deis/deis/commit/9080b7766b2bf326fd2f446e1e5fd5c0cc2588a7) releases: update procedure for Deis releases
 - [`17e2d1f`](https://github.com/deis/deis/commit/17e2d1f54fc34d0af86f81d0f59207b9619c60cb) managing_deis: rewrite several pages based on store
 - [`fb7d9c4`](https://github.com/deis/deis/commit/fb7d9c40b12e44092d90f16f881ce7419036ed92) client: update docs with preferred deis CLI install method
 - [`009ddf7`](https://github.com/deis/deis/commit/009ddf7685269dbdd036831128ab851fcbaefedc) (all): changed etcdctl references to deisctl
 - [`779fdad`](https://github.com/deis/deis/commit/779fdad57260c2ab3c676a50ff958eebd5e37d16) managing_deis: add store documentation
 - [`ca16746`](https://github.com/deis/deis/commit/ca167467dfcc7ada68ca2e77ca7e19e787a4fdc6) deisctl: remove sudo from deisctl installer examples
 - [`104e6a8`](https://github.com/deis/deis/commit/104e6a8036016ee33e2be748c4b5038c3dd4b56e) standards: remove [skip ci] instructions
 - [`e5d5c1e`](https://github.com/deis/deis/commit/e5d5c1e020a6a5ca9ecdb289a838594541dd2279) client readme: remove travis build status icon, as travis is no longer being used
 - [`cbb3b0f`](https://github.com/deis/deis/commit/cbb3b0f038fe09548c67771e7f4e40a8e25387a5) logspout: update docs to reflect deis' fork
 - [`188e30a`](https://github.com/deis/deis/commit/188e30a6b852b10d4225d3fd3705d051763373dd) deisctl: document DEISCTL_UNITS better
 - [`ded85ec`](https://github.com/deis/deis/commit/ded85ece8995538fd27d353fa1050bfbec4b656c) installing_deis: remove xip.io instructions for ELB
 - [`4b81149`](https://github.com/deis/deis/commit/4b81149cdfe57f490bc5590c6f049659accbdcf5) deisctl: clarify deisctl version when using install.sh
 - [`6f096af`](https://github.com/deis/deis/commit/6f096af4c85a1625cf317bb319627d6537cf1975) vagrant: bump version to v1.6.5
 - [`ed1dba5`](https://github.com/deis/deis/commit/ed1dba5fa37927d6494a630c72cbfc6f9988a93d) contrib/rackspace: add manual update instructions"
 - [`5550a3a`](https://github.com/deis/deis/commit/5550a3afc683afc8ff669adf3e4a7f68348ae916) deisctl: add unit search paths behavior to README
 - [`bafe331`](https://github.com/deis/deis/commit/bafe331f8c2f972946a1554c35ba612b514d919e) CHANGELOG: update for v0.12.0
 - [`1c30515`](https://github.com/deis/deis/commit/1c30515f5aff816167159131ab5fc0cec0c31557) README: update dev documentation
 - [`97b553a`](https://github.com/deis/deis/commit/97b553af90a24a3c4494ac3d16a9118fafa9d966) README: link to latest installers on S3, omit "how to build"
 - [`486b422`](https://github.com/deis/deis/commit/486b4225db58f1bfd44437e094179c1ec5c7280e) readme: update install instructions
 - [`e65ffcb`](https://github.com/deis/deis/commit/e65ffcbf1968ce1acf7105c938c86d2a40971d0d) readme: minor language updates
 - [`b08541d`](https://github.com/deis/deis/commit/b08541dd6da1d32c57b153e217f41df8c919ac2e) readme: first pass at readme

#### Maintenance

 - [`4eb0089`](https://github.com/deis/deis/commit/4eb00899ea5e58a1cab2a72f0fe54e2fc5d9870e) registry: use store component for filesystem layers
 - [`893414a`](https://github.com/deis/deis/commit/893414a5f0216851caa353efc4897d95e5cc9c12) (all): remove deis-database-data and deis-registry-data
 - [`b5cb742`](https://github.com/deis/deis/commit/b5cb7427b3e83381afa4aa6e52c29231c5a18862) store: add deis-store component
 - [`4bdfcb3`](https://github.com/deis/deis/commit/4bdfcb317e2054a0b299243b8e11353ffe50c4f9) contrib/gce: bump CoreOS to 459.0.0
 - [`f83ceae`](https://github.com/deis/deis/commit/f83ceae687340f71e5af98580260851809002a64) (all): bump CoreOS to 459.0.0
 - [`7a1feec`](https://github.com/deis/deis/commit/7a1feece3669ef85854f05e56949b2ae3c4754a6) controller: remove unused django-yamlfield from requirements
 - [`5756746`](https://github.com/deis/deis/commit/57567465c63479765b988faf61b390c0e40fdc5d) controller: bump gunicorn to 19.1.1
 - [`b9fe218`](https://github.com/deis/deis/commit/b9fe21830b1d986bfe1905721e7780f18a9abb8f) builder: migrate to cedar-14 stack
 - [`70969eb`](https://github.com/deis/deis/commit/70969eb56a8ee97b9d77971e2cf071e6d7a49ceb) pip: update pip installs to 1.5.6
 - [`c5e6e01`](https://github.com/deis/deis/commit/c5e6e01c98192aa5b5143c44029898d98dd449e4) (all): bump CoreOS to 452.0.0
 - [`d45dc77`](https://github.com/deis/deis/commit/d45dc77dba02ba5d4df186c27ac653292c3dacf8) tests: add log_phase helper
 - [`420bc97`](https://github.com/deis/deis/commit/420bc97afcafddc58abbc7e291ef1e3c3b547043) (all): Rename PUBLISH to EXTERNAL_PORT
 - [`7595033`](https://github.com/deis/deis/commit/75950332cade70fba1602f3513881af5990430f0) cache: Remove moved systemd unit file
 - [`9faf9fa`](https://github.com/deis/deis/commit/9faf9fac7bdbf03c4ef8a2530812b777576071de) (all): Rename PUBLISH to EXTERNAL_PORT
 - [`ddb34d3`](https://github.com/deis/deis/commit/ddb34d3bc309cf2f0a402bcb9112a785f5bbba9e) (all): remove deis-builder-data
 - [`ccb48ae`](https://github.com/deis/deis/commit/ccb48aee7585fb5869168fc77822d3002b25b641) (all): remove deis-builder-data
 - [`05eeb77`](https://github.com/deis/deis/commit/05eeb771604a123bfe9cd3d8d55a422f37a9ca4f) deis: bump version to 0.13.0-dev
 - [`f4446eb`](https://github.com/deis/deis/commit/f4446eb2ffada63e7513d84f4d84360eb7e153ef) deisctl: bump version to v0.13.0-dev
 - [`2b93db0`](https://github.com/deis/deis/commit/2b93db06d0d5e79fa00e834935f5978960841e05) release: update version to v0.12.0
 - [`1a38aff`](https://github.com/deis/deis/commit/1a38aff5d22a7bee965029e54fbdeb4b3a0b8fbb) units: remove deprecated X-Condition from fleet units
 - [`006556e`](https://github.com/deis/deis/commit/006556e9ef8c437d0a1204b00c6e4902dd574e92) README: update current version to 0.12.0-dev
 - [`539ed23`](https://github.com/deis/deis/commit/539ed23e7d972b6f8bda81b92f618ed09f8c08a3) deictl: bump version in sync with Deis
 - [`23301e8`](https://github.com/deis/deis/commit/23301e84b52fb0354cbcfcad451f96627a8c4d53) godeps: bump fleet, updateservicectl, docker
 - [`3d0bf7f`](https://github.com/deis/deis/commit/3d0bf7f17d916cf24f15cde34da4707a519862a8) flags: switch to DEISCTL_TUNNEL
 - [`b610417`](https://github.com/deis/deis/commit/b6104173b1c6fadcf311650335b556670fcdbcc4) version: bump to 0.11.0

### v0.11.0 -> v0.12.0

#### Features

 - [`5cf53e9`](https://github.com/deis/deis/commit/5cf53e90ebfc5cd328da9289c88a540b8956406a) tests: add test-nightly.sh script for CI
 - [`c3b262a`](https://github.com/deis/deis/commit/c3b262a053f853910aab4d4652b485fbbba1f9bc) tests: add end-to-end acceptance test script for vagrant
 - [`16e92ee`](https://github.com/deis/deis/commit/16e92ee919ae9f8414a270a6c2572bec719cbb74) controller: inject tag version value in environment.
 - [`49daee2`](https://github.com/deis/deis/commit/49daee2a86cad07f558700a38d10b0b74ea551c0) controller: inject tag version value in environment.
 - [`940752b`](https://github.com/deis/deis/commit/940752bb8b09749f90f250fb637597f3c7ca7bce) version: add version constant
 - [`32b3aa5`](https://github.com/deis/deis/commit/32b3aa54cab2ab74b2084d7fb92b1e814e821742) version: add -dev flag for unversioned releases
 - [`f9a736b`](https://github.com/deis/deis/commit/f9a736b5316876229878c05c92380ac7837a5369) router: custom timeouts for builder and controller
 - [`e5488e8`](https://github.com/deis/deis/commit/e5488e892dbe20fe5e698e43be13de61a5265edb) router: custom request size.
 - [`9b72b74`](https://github.com/deis/deis/commit/9b72b74d7ca178adfc251f17b94919c28c574b5a) client: add colorized logging
 - [`9bc583c`](https://github.com/deis/deis/commit/9bc583cfe94a36514a701f48410c7d55102ae0d2) tests: add CI node setup script
 - [`4ed700b`](https://github.com/deis/deis/commit/4ed700bd07cc740950857fcbabfaa09b0c94052e) controller: pipe release notes to app logs
 - [`c1c1017`](https://github.com/deis/deis/commit/c1c10172c953018a616e6a79f8526239ea7abaf5) contrib: Blacklist DigitalOcean regions.

#### Fixes

 - [`d5b419e`](https://github.com/deis/deis/commit/d5b419e055cf3446569a0892bfc979120443e04c) tests: make test-nightly.sh executable
 - [`7aa8062`](https://github.com/deis/deis/commit/7aa8062d46509379032b4c136e4d4f73e38076fe) controller: fix scale log event on git push
 - [`d6f7dd0`](https://github.com/deis/deis/commit/d6f7dd0548255d69516730db74fb7273220ef2b9) contrib/gce: fix docker storage being wiped on boot
 - [`2a8488c`](https://github.com/deis/deis/commit/2a8488ccdafb29a26a0a79f5bec2028e4811fdb4) Vagrantfile: specify AMD network to work around transfer issues
 - [`2220271`](https://github.com/deis/deis/commit/22202712e863c7888b83d61ffcdfde5f70767a55) tests: improve functional test setup logic
 - [`5bc1c80`](https://github.com/deis/deis/commit/5bc1c809b90cf0fe3e8c09b11261aa19b0f789ce) user-data: add fixed version of sed command to "make discovery-url"
 - [`27860be`](https://github.com/deis/deis/commit/27860bee20e38a3a9eceea464d383253ecd6af0c) makefile: default to 3 nodes, test for docker binary
 - [`0943823`](https://github.com/deis/deis/commit/0943823d28fd5c397e14951f7ada83c8bbd4828d) readme: address feedback, fix typos
 - [`8d08d43`](https://github.com/deis/deis/commit/8d08d430d0b65cfc12e5983b88d83c9a5a114105) user-data: restore default discovery url
 - [`d34b839`](https://github.com/deis/deis/commit/d34b8393e7a760cdb8c68126f08b280a708f7fa2) controller: properly purge domain etcd keys on app deleteion
 - [`0920b9e`](https://github.com/deis/deis/commit/0920b9ee039f8b3ad5d1db8b0fb9cd33f117252d) client: also encode keys which could be int
 - [`66333da`](https://github.com/deis/deis/commit/66333da4ab1797d5836811fb3b26841dfe4637ae) client: fix syntax error in config:list
 - [`79dd499`](https://github.com/deis/deis/commit/79dd4991fee82ad3d1b633bc7d9830592312d605) contrib/ec2: use conditional root DeviceName
 - [`99f0152`](https://github.com/deis/deis/commit/99f01525d9fe830055385876f1a35dc5bd5bcca9) controller+client: utf-8 encode only string types
 - [`8092009`](https://github.com/deis/deis/commit/8092009b40d55974fdf13e9fdb9cb2739df729aa) tests: temporary workaround to try and get jenkins passing
 - [`e3780fd`](https://github.com/deis/deis/commit/e3780fd5f8d5a1a12600a43d532e42226ccf66d9) client: use only stdout handler
 - [`8b0158e`](https://github.com/deis/deis/commit/8b0158ec2a8b4e363399e4cc66002489b91c53c9) scheduler: add gunicorn timeout for long-running methods
 - [`0cf4886`](https://github.com/deis/deis/commit/0cf488677aafc91a8828bbf611db2303e52ea434) scheduler: handle announcer timeouts to work around https://github.com/docker/docker/issues/8022
 - [`b8c9a43`](https://github.com/deis/deis/commit/b8c9a430b0dce59eb38d9a1ee603c99b9188c0a0) contrib/coreos: bump etcd peer heartbeat to 500ms
 - [`b2b8cbc`](https://github.com/deis/deis/commit/b2b8cbc5c4d168c814b075d60171aab02ef0ebbe) logger: properly parse incoming syslog messages"
 - [`230a8c6`](https://github.com/deis/deis/commit/230a8c685f8f9fefe335d1167eb008443a5bb586) client: convert client.yaml to client.json once if needed
 - [`1d20351`](https://github.com/deis/deis/commit/1d203512f9a38986679a7a91681697478a7fd223) tests: use our own fork of postgresql container for testing
 - [`362b042`](https://github.com/deis/deis/commit/362b04207430679bd3c2c57ba6ee37d4c2be4cae) builder: hardcode STACK to cedar
 - [`c15d909`](https://github.com/deis/deis/commit/c15d909c4dfa58f28b82c6f0e07b7cc730606ac9) controller: allow unicode in config:set and app logging
 - [`7df6b45`](https://github.com/deis/deis/commit/7df6b45ac384ce6b3c98a1c28ed6aa1c66dd2ab6) client: concatenate strings together
 - [`4f78bfd`](https://github.com/deis/deis/commit/4f78bfd788dc7e67d5a4537ae3058a02c8232e51) docs: update links to buildpacks included with Deis
 - [`7c9fce2`](https://github.com/deis/deis/commit/7c9fce2d83cec93cd00a1e3247e838e27bd0c49a) test: Verify apps:run output
 - [`dd0aeb0`](https://github.com/deis/deis/commit/dd0aeb09b7a9242fa53b901d00836224772a4be4) client: do not prepend the app name twice
 - [`0df66bd`](https://github.com/deis/deis/commit/0df66bdd966f5d3fd9b02825b7194c9ce891cd1d) logger: remove priority from log messages
 - [`0c66ecd`](https://github.com/deis/deis/commit/0c66ecd629a6b41d95fd80e149b52e724ba9405d) logger: properly parse incoming syslog messages
 - [`3b7edd7`](https://github.com/deis/deis/commit/3b7edd7e0fffc11b24a55a51c3373e31cef4c036) tests: call other tests scripts in local dir by explicit path
 - [`c273b27`](https://github.com/deis/deis/commit/c273b27573e3ba0c209e3f1c9ba6ee9af3463f82) scheduler: initially wait for container to start
 - [`56c6810`](https://github.com/deis/deis/commit/56c68103f7d663882c6645516ec315e604a43d96) contrib/digitalocean: update network interface to eth1
 - [`ce220ab`](https://github.com/deis/deis/commit/ce220abb1957b3361cc3c2a1a9cb6c614779adec) contrib/coreos: bump timeouts for etcd and fleet
 - [`42602c8`](https://github.com/deis/deis/commit/42602c802276623232a36fb27d1c904b172e9609) contrib/gce: fix accidental typo in GCE generator

#### Documentation

 - [`081e289`](https://github.com/deis/deis/commit/081e2899484cd0c1ef84b25170e7553f6ad3dfc1) CHANGELOG: update for v0.12.0
 - [`67d3a45`](https://github.com/deis/deis/commit/67d3a4589a60d0dad7d7d460bd3bb3245e92511c) releases: update release procedure with deisctl tasks
 - [`a8823bc`](https://github.com/deis/deis/commit/a8823bcc9707a1554bea0f2d2e63fe76c055c841) contrib: update docs with deisctl
 - [`4e9fcb1`](https://github.com/deis/deis/commit/4e9fcb1bb375d7bbb2f14df39cca169e0a1e00dd) using_deis: add Using Docker Images section
 - [`8f70ad9`](https://github.com/deis/deis/commit/8f70ad98a49553683091d6380f0b9f4f28a25037) readme: cleanup readme after deisctl integration
 - [`a40b17b`](https://github.com/deis/deis/commit/a40b17b08e8300a9f8dd34a73f67da2b369b5ecf) (all): update documentation for new deisctl provisioning workflow
 - [`466b689`](https://github.com/deis/deis/commit/466b68931328ec056cf06f9fd653269b3ccbeaad) readme: add hack instructions to readme, remove newline
 - [`f2e8117`](https://github.com/deis/deis/commit/f2e81177ce078bc8c0aaa0fb41b146ae3388af72) hacking: new docs for hacking on deis
 - [`9267de1`](https://github.com/deis/deis/commit/9267de13c0035717148551989f4cc825a68b6dba) readme: update root readme with new deisctl instructions
 - [`0a70a3f`](https://github.com/deis/deis/commit/0a70a3f4e1c80ad7a050be46d6aa46810339c6bf) installing_deis: remove in-place upgrade documentation
 - [`911ba8e`](https://github.com/deis/deis/commit/911ba8e1ce179b14654b740c7dbd667cdaf011df) router: keep etcd keys ordered.
 - [`ee83813`](https://github.com/deis/deis/commit/ee83813ba4372f6f3b7be51295a885c8cfc061bc) EC2: CloudFormation
 - [`d26c641`](https://github.com/deis/deis/commit/d26c641ecffcf5f8b79b68c5aa342bf7e552ef82) ssl_endpoints: open port 443 in ec2 secgroup
 - [`b6ccd8b`](https://github.com/deis/deis/commit/b6ccd8b65470237c03f52ae542348de8ba9fbbb5) client: update with dev CLI binaries and local deis.py suggestion
 - [`8523d59`](https://github.com/deis/deis/commit/8523d5930d4efdc9e649527e527858c7531d9e70) installing: deploy new cluster with external components

#### Maintenance

 - [`f95204b`](https://github.com/deis/deis/commit/f95204bf2d3439faac53fc0f75730f7de038b923) release: update version to v0.12.0
 - [`e70af2f`](https://github.com/deis/deis/commit/e70af2f6ee259052e41ab17996a7ce3a11139220) router: Updated nginx to 1.6.2
 - [`fc5ea45`](https://github.com/deis/deis/commit/fc5ea4516b2c25679fdf2d194f1455d2c6b92021) (all): remove fleet/systemd units
 - [`8519f79`](https://github.com/deis/deis/commit/8519f795eae28fdf8010e335641007e3c878519c) (all): bump to CoreOS 444.0.0; fleet 0.8.1
 - [`d014537`](https://github.com/deis/deis/commit/d0145376cf37e58bdcbebb347871657da91c1652) utils: move encode to api/utils module
 - [`050fb06`](https://github.com/deis/deis/commit/050fb061e7c2acfd49605a574839c12e4d6a20c2) (all)/systemd: remove deprecated X-Condition in fleet units
 - [`d1fcdd5`](https://github.com/deis/deis/commit/d1fcdd55c45b4d7d9a5620334fbd3ef7bd318a0d) (all): update CoreOS to 440.0.0
 - [`a97cdc3`](https://github.com/deis/deis/commit/a97cdc3a59a375df46ea8d63dc45461ef430a919) (all): bump CoreOS to 438.0.0; Docker to 1.2.0; fleet to 0.8.0
 - [`a23719e`](https://github.com/deis/deis/commit/a23719ee39175e185577214cebc072fba06767ef) registry: bump version to v0.8.1
 - [`d40d2bb`](https://github.com/deis/deis/commit/d40d2bbe7eb6aa359141924192536ea48234618c) controller: update Django to 1.6.7 bugfix release
 - [`9a692ea`](https://github.com/deis/deis/commit/9a692ea77ecdf9409fef804dde7b969d61149ca3) builder: pin cedarish to cedar stack
 - [`7c518db`](https://github.com/deis/deis/commit/7c518db4a6c0373c2829a5a187d852c7990ac7d5) controller: update django-guardian to 1.2.4
 - [`a4faa55`](https://github.com/deis/deis/commit/a4faa5560e96583fe94a4b9841bdfb627306be19) (all): rename deprecated X-ConditionMachineBootID
 - [`8c72985`](https://github.com/deis/deis/commit/8c72985a8f0894a5c4830ee93a4be382ec55ac85) builder: remove unnecessary slug* documents
 - [`cc339d8`](https://github.com/deis/deis/commit/cc339d851b24bbec01f99517227b79642a0a6ba2) builder: import deis/slugrunner
 - [`30228c9`](https://github.com/deis/deis/commit/30228c9fd3dc97be7e4f2bc11bd71947a9b83bee) builder: import deis/slugbuilder
 - [`db8df5e`](https://github.com/deis/deis/commit/db8df5eb94028a9f9d717fe3da588e7c0a39d72d) controller: update django to 1.6.6 security release
 - [`21ddb96`](https://github.com/deis/deis/commit/21ddb9680e984be4e42c2835fe4d1a422cd0e73c) (all): update master to v0.12.0

### v0.10.0 -> v0.11.0

#### Features

 - [`df47b06`](https://github.com/deis/deis/commit/df47b06cf7d68be3d2a0398ec447066818ae030e) (all): add `make test-style` targets
 - [`59a5796`](https://github.com/deis/deis/commit/59a5796085613f720df731fdc1b5f33ebc7e0278) tests: integration tests for tags, move limit tests under config
 - [`5b0bbec`](https://github.com/deis/deis/commit/5b0bbecc4f459dc7fa9502b8bbd050096bc11748) client: support for application key/value tags
 - [`517869d`](https://github.com/deis/deis/commit/517869dd9ba53939ac741ad9b92ff81efb85d0b3) controller: support for application key/value tags
 - [`19396a1`](https://github.com/deis/deis/commit/19396a1350d7cb953bcb5d813c6c556261c94f68) client: add proxy support
 - [`136e1fb`](https://github.com/deis/deis/commit/136e1fbe8ff399bc0b61ce2372a516c8fc32e539) slugbuilder: build without using a pipe for git archive
 - [`8bea984`](https://github.com/deis/deis/commit/8bea9845d49430d4b6ca64d3e8d30f3c1361a6a0) builder: automatically run `git gc` on deploy
 - [`8c335e1`](https://github.com/deis/deis/commit/8c335e18854b0eae94cbfae045360e8c8295544e) controller: expose app url
 - [`4da2426`](https://github.com/deis/deis/commit/4da2426a4575f4c8e2c434a5579889b560356c3a) contrib/ec2: always deploy into VPC; support all instance types
 - [`4003921`](https://github.com/deis/deis/commit/4003921dba4664f5498465b9201538b2cc48ea4f) client: add commands for managing memory/cpu limits
 - [`b01033b`](https://github.com/deis/deis/commit/b01033b4e8d98f860d569f7c5fae1fe3369a1efe) controller: add endpoint and infrastructure for limiting memory/cpu
 - [`40b2929`](https://github.com/deis/deis/commit/40b2929ea1ab3d8153a61d957e648dcd090154b2) controller: allow tags for `deis pull`
 - [`309c6ce`](https://github.com/deis/deis/commit/309c6ce2cf53985a538363269f0e134dc6f5005b) router: Support for WebSockets
 - [`b0572ef`](https://github.com/deis/deis/commit/b0572ef9a987b2b8c9d02a2943b1679ab2cd96cf) builder: allow custom slugbuilder and slugrunner images
 - [`3e4e44b`](https://github.com/deis/deis/commit/3e4e44bd6553b1e5fc47639c7422da7ba5fcb3f6) client: allow cluster id rename
 - [`9e645f6`](https://github.com/deis/deis/commit/9e645f616857e51c5dfcdd5ba7963140bf23dda9) client: add -a option to client
 - [`e004275`](https://github.com/deis/deis/commit/e004275102c73dfba0465affe8d3a6cd7dc9ad51) contrib/openstack: script for deploying deis on OpenStack
 - [`99dfdf1`](https://github.com/deis/deis/commit/99dfdf18ebe510db668d317be53d269a7efe0ede) client: add config:pull command
 - [`45fb9a6`](https://github.com/deis/deis/commit/45fb9a6f0d3265511927b3cf3f6dad593d90dee5) registry: Configurable storage backends
 - [`d5102ed`](https://github.com/deis/deis/commit/d5102ed37c980d72575b7731a462bedea11e57d3) builder: cache progrium/cedarish image
 - [`e079ffa`](https://github.com/deis/deis/commit/e079ffa5c7b2e258b732a111cbd76efc4c2317e9) contrib: allow private networked regions
 - [`975fb1a`](https://github.com/deis/deis/commit/975fb1ace99250edae4b0a960987079962b3562e) contrib: add usage to DigitalOcean script
 - [`f9d47e3`](https://github.com/deis/deis/commit/f9d47e3f77cf93ed28ef3d135804be47ab609525) client: check server version
 - [`22b3ba6`](https://github.com/deis/deis/commit/22b3ba653be8f71845257759f61d6c739751d209) router: Avoid nginx 502, 503, 504 errors

#### Fixes

 - [`4753474`](https://github.com/deis/deis/commit/47534748ad6479a4cf54af6454faee71c4871709) docs: add client reference for `deis tags`
 - [`1ff20cc`](https://github.com/deis/deis/commit/1ff20cc390ae28d3f1f047305ad96d4ff7b4f62c) tests: take a nap after fleet restart
 - [`a15bbaf`](https://github.com/deis/deis/commit/a15bbaff9f4b23e5f6c28d0e83c0691564ef3f4a) router: add check command for nginx config
 - [`e7a844f`](https://github.com/deis/deis/commit/e7a844ffa56afe3baf8ecc7089edc133259aef50) (all): all services ensure valid config file on reload
 - [`adf6cc1`](https://github.com/deis/deis/commit/adf6cc1090224fd51ce2fa0192bd393c714f0875) controller: validate Cluster.domain field
 - [`1a53a73`](https://github.com/deis/deis/commit/1a53a73919bf4c578bb83b29c40314ba522102b6) controller: validate app structure values
 - [`a1c1432`](https://github.com/deis/deis/commit/a1c1432be071b8eae68f870546c26b837c5edc20) controller: validate Cluster.hosts field as comma-separated hostnames
 - [`af2ffcf`](https://github.com/deis/deis/commit/af2ffcf0b28c1176da47a2bae5ae782e1e8def7c) builder: don't reload cedarish if its already loaded
 - [`66056e6`](https://github.com/deis/deis/commit/66056e6f3f64d9f4253209b37f72e88b74567c27) builder: disallow password authentication
 - [`0a518d6`](https://github.com/deis/deis/commit/0a518d6941185136e9d2b608e12ea2fb816a6e3a) controller: specify dicts as defaults and guard against old errors
 - [`4a0ad8d`](https://github.com/deis/deis/commit/4a0ad8d572d57a4b8881c70983ee080fc50cdfd2) controller: store limit fields on config
 - [`007240c`](https://github.com/deis/deis/commit/007240c481a447b6b1ae9cae20834c44bed6536b) docs: silence header warnings and restore broken links
 - [`78b0f04`](https://github.com/deis/deis/commit/78b0f04087c6b58c8ee1ac075ecd7a2d3c0ab01d) controller: don't encode empty JSON objects or arrays as strings
 - [`b0b4b5a`](https://github.com/deis/deis/commit/b0b4b5a62fc39f3b1a7819188949ec954695134d) client: set cookies.txt readable only by current user
 - [`b44a58f`](https://github.com/deis/deis/commit/b44a58ffaf2891d83b2d0f127630c20af58f98cf) controller: use /bin/sh as entrypoint for run"
 - [`d4b284f`](https://github.com/deis/deis/commit/d4b284f6564c99cbbee922e1f9109ca5ff0cdab7) docs: remove coveralls badge and mentions of coverage stats
 - [`423a55d`](https://github.com/deis/deis/commit/423a55d95bd6eb54171cc0f1f599b65f9b540862) controller: honor app permissions correctly on app-related views
 - [`e5984fe`](https://github.com/deis/deis/commit/e5984fe87f063010e646d2ccc6cfba2eb3b2169b) builder: set /app as working directory for apps
 - [`f8168c3`](https://github.com/deis/deis/commit/f8168c301f055c891870b22f18e24fa0134ab224) controller: set celery log level from envvar
 - [`a882f1f`](https://github.com/deis/deis/commit/a882f1fff70d768a86c71d215fde3802bd3abe64) controller: check image for hostname
 - [`987b2d5`](https://github.com/deis/deis/commit/987b2d56109c9442ad3febc2dc1ee89482505dde) controller: chown log mount to deis user
 - [`a50158b`](https://github.com/deis/deis/commit/a50158be89a5cdd5a1c45de89b9a2a5cff75f8b0) controller: remove application logs on delete
 - [`79e2b75`](https://github.com/deis/deis/commit/79e2b75984ded1542a445f7880c6bd6f189f1175) controller: fix logging of limit changes on release
 - [`29b6138`](https://github.com/deis/deis/commit/29b6138b76eb92f3d3653500c422d92e435680a3) builder: fix race condition in /bin/boot
 - [`93eb76b`](https://github.com/deis/deis/commit/93eb76b46b60221d5f8ecb090b701a23873a08ac) docs: correct the upgrade flag syntax for pip
 - [`f2259cf`](https://github.com/deis/deis/commit/f2259cf45edf94f65fb0c497812e54ec993f5a64) Makefile: fix check-fleet Makefile target
 - [`9028697`](https://github.com/deis/deis/commit/902869746deef943ed8b2bf5d8fcd9c3e8c1a6fe) tests: remove GetCommand(
 - [`3339045`](https://github.com/deis/deis/commit/3339045dd506d0fec136239e735d2153977f23c2) controller: rename /limit to /limits
 - [`df7fb05`](https://github.com/deis/deis/commit/df7fb0563d03abe5d3e0d29803054b7f6af74096) client: move options list before arguments
 - [`eb12271`](https://github.com/deis/deis/commit/eb1227104762a589af3d6f45b4218e61f3173e22) builder: bump timeout to 30m
 - [`c5922a9`](https://github.com/deis/deis/commit/c5922a9aba4f75a1edf9b5421ed78b22f8ec070d) controller: tack on host/port to image id
 - [`8bb9bcd`](https://github.com/deis/deis/commit/8bb9bcd27eab05857263416e59574af41616c43f) controller: remove /run/deis/determine_registry
 - [`e9082d4`](https://github.com/deis/deis/commit/e9082d4dc4258c4732af436b7af6b69113ab9b4d) controller: Determine registry address during start
 - [`b30c95c`](https://github.com/deis/deis/commit/b30c95cb19852964976d1abe3e979ace726c8a94) contrib/digitalocean: fix region check
 - [`b071e08`](https://github.com/deis/deis/commit/b071e088c86cf296d9045aaf9626bf464480ece4) tests: allow env vars to override important test settings
 - [`16f33bf`](https://github.com/deis/deis/commit/16f33bf8906d17b6aaa21d346b33de155ece112d) controller/registry: strip hostname from repo
 - [`aaa95b1`](https://github.com/deis/deis/commit/aaa95b1308450081b9d9d3b5a7e81958beaf5980) makefile: delete missing files on rsync
 - [`1219340`](https://github.com/deis/deis/commit/121934003e028b53eb0e29061bc82225e2ff15ac) contrib: read proper argument
 - [`a82b54c`](https://github.com/deis/deis/commit/a82b54c4c0bc9b964a92136409777526ad5217f9) contrib: remove unnecessary call to echo
 - [`da8e99a`](https://github.com/deis/deis/commit/da8e99aab04ba4421bf88550ed2aa36312386849) controller: add south_migrations to flake8
 - [`727d400`](https://github.com/deis/deis/commit/727d4008d27b144c46bb35d93e3f0314c4e0e370) contrib/coreos: configure public IP for fleet
 - [`47a44f4`](https://github.com/deis/deis/commit/47a44f4a4b0a433a57d63ecd66a566a0495cb64a) Makefile: submit only *.service units for registry and controller
 - [`8c9a8ac`](https://github.com/deis/deis/commit/8c9a8ac6e83ea899604e877381f628e09350866c) client: make --auth optional
 - [`3e367f3`](https://github.com/deis/deis/commit/3e367f3355426782ac98fa6cb2f31726645d6569) client: remove check for error message
 - [`1df8c99`](https://github.com/deis/deis/commit/1df8c99ef22ee7d89b991ab1ce1c5fd6b78e0fc5) builder: pass in http_proxy env variables from host
 - [`d3bd219`](https://github.com/deis/deis/commit/d3bd219c24a360b30b978859c477881466153ce9) client: skip logout if connection fails
 - [`32fd600`](https://github.com/deis/deis/commit/32fd60003d474d4c7d82e2992904dd26e3f12b59) Makefile: force constant output of fleetctl list-units
 - [`87c6bf5`](https://github.com/deis/deis/commit/87c6bf5a0be17822eb626dfd3924cfcd69f9975c) userdata template: put `nse` function in FHS
 - [`cae7a8f`](https://github.com/deis/deis/commit/cae7a8f771ef10bb93378ebf7fe418b058e5bc18) userdata template: put motd directly in the OS
 - [`fbe7efa`](https://github.com/deis/deis/commit/fbe7efa9a7c46dfe448a7988751a396adb95093e) tests: use aufs in builder functional tests
 - [`d47f790`](https://github.com/deis/deis/commit/d47f790c70e3bb32b3dd7e28a29f9ac34238e464) tests: report curl output if it wasn't "Powered by Deis"
 - [`cb07b96`](https://github.com/deis/deis/commit/cb07b96a6285b7dbeb7abc6d8faf2f2c18738e6d) client: restore `--cluster=dev` as apps:create default
 - [`06aa80e`](https://github.com/deis/deis/commit/06aa80e57113099fa5118e1f3285d409084a8f0f) tests: don't make assumptions about the user's system
 - [`98edce4`](https://github.com/deis/deis/commit/98edce408ba761cd5919654eda754e2a6a95130b) controller: disallow uppercase app URLs
 - [`3f35905`](https://github.com/deis/deis/commit/3f35905b12ddc0680d843c643e90dfb9f1a8cb26) docs: improper config:pull docstring syntax
 - [`2f3dd93`](https://github.com/deis/deis/commit/2f3dd93288cd22feba92ea61a37cc960bede457b) contrib: give /var/lib/docker 30GB on DO
 - [`f82f5ad`](https://github.com/deis/deis/commit/f82f5add10f78970913392381f4e11ddfb34814c) (all): etcdctl and docker history commands eat error output
 - [`8e01b9d`](https://github.com/deis/deis/commit/8e01b9d46b6e751d7f31284532d7ad24902ad055) client: `deis run` no longer requires remote
 - [`6bb034c`](https://github.com/deis/deis/commit/6bb034c364dac600d0aa4dfcf16698fc28a67f52) controller: reverse version output
 - [`4742e74`](https://github.com/deis/deis/commit/4742e7494ffdc0f4417b353f9e6149b6b3521257) Makefile: remove redundant command parameter
 - [`6cc050c`](https://github.com/deis/deis/commit/6cc050c72e153e6764b3ce43e2d7c0f02e5fa5b2) client: KeyError on cookie retrieval
 - [`c80d3dd`](https://github.com/deis/deis/commit/c80d3dd8f4a6489255eb8a7d14bd2a710b89f461) fleetctl.sh: clean up unitfiles from the filesystem
 - [`5703672`](https://github.com/deis/deis/commit/5703672fdfba059f5cf2e1892955e75910ec78f2) client:  add sshkey encode in client update api call

#### Documentation

 - [`6e70e27`](https://github.com/deis/deis/commit/6e70e27fae9208fe1ebd3ff6d1df5bc881f89fa6) using_deis: clean up proxy support
 - [`e6d3452`](https://github.com/deis/deis/commit/e6d345286cd60e69c06367fb4c1f74d591a56d3d) community: add policy on bounties
 - [`88a23aa`](https://github.com/deis/deis/commit/88a23aaa41f3d37930238aeabba488ef6d19ee6f) installing: fix typo in load balancers doc
 - [`e04db78`](https://github.com/deis/deis/commit/e04db78060bd008c6477c6977cd6f45a927e70d4) contrib/ec2: document recommended instance size
 - [`fb983a2`](https://github.com/deis/deis/commit/fb983a2c7e06648759b77cf3a9e6c45ffba41604) client: specify default env for `deis run`
 - [`5460f73`](https://github.com/deis/deis/commit/5460f735259a2ed8c902746ed0908140ee3c391c) reference: add `deis pull` reference
 - [`61e95d5`](https://github.com/deis/deis/commit/61e95d5c585b792fa26be6be88f5ef6df93a23bd) contributing/localdev: update localdev docs
 - [`ed22a9c`](https://github.com/deis/deis/commit/ed22a9c1e2873ace656a1168c116b510b364b838) installing_deis: clarify no web UI
 - [`11699fa`](https://github.com/deis/deis/commit/11699fa22b646e2773a1e67be7cb1c4efaa6eb4c) gce: instructions for running Deis in Google Compute Engine
 - [`d4df11e`](https://github.com/deis/deis/commit/d4df11ef00c13dfb6067196a0eec84d56eaad406) dockerfiles: only allow one port exposed
 - [`25b4217`](https://github.com/deis/deis/commit/25b42179c45af2b92805a89a98fc83043183b277) (all): clarify that auth SSH key for clusters cannot have a password
 - [`08a4140`](https://github.com/deis/deis/commit/08a41400d068e343a9dab21a1942ca8c68972d7e) README.md: clarify current release and encourage use of master
 - [`46513d6`](https://github.com/deis/deis/commit/46513d6b809236a123e5e70c14b22a485710cb0a) managing_deis: add custom slugrunner and slugbuilder to docs
 - [`8754a6f`](https://github.com/deis/deis/commit/8754a6fb75c333f95abb398ebbfb73cf865d1660) logger/syslog: amend syslog server docs
 - [`247e521`](https://github.com/deis/deis/commit/247e521cd521f8678817dd990bf21a43aeb97cfe) contrib: expand bare-metal provisioning
 - [`b099e71`](https://github.com/deis/deis/commit/b099e71390d515a58139f8fdc4ed97ff92f82447) (all): remove references to `DOCKER_HOST`
 - [`f46b7ea`](https://github.com/deis/deis/commit/f46b7ea17e33f38d32cc8cfa732d010591853db6) contrib/rackspace: add manual update instructions
 - [`898bb61`](https://github.com/deis/deis/commit/898bb613b81e3c343e25b954c23a13cde280b933) contrib: add REGION_ID as optional argument
 - [`1d40560`](https://github.com/deis/deis/commit/1d40560a1696df234823c9b3b5e7f77a63e1d390) controller: add Limit server autodocs
 - [`764522d`](https://github.com/deis/deis/commit/764522d88159f75e7341c534bdae808b5ec5be6d) client: more clarity on limits:set
 - [`773f105`](https://github.com/deis/deis/commit/773f105febf4b279ee3fda494848b49a676a7a1f) reference: add deis limit autodocs
 - [`c42acb9`](https://github.com/deis/deis/commit/c42acb97c6d41d760411bd70a502c2ad17a08ba8) client: more clarity on limit:set options
 - [`a0ee3bd`](https://github.com/deis/deis/commit/a0ee3bd86cca2bda72179eb45c99d89366178bcf) client: add more clarity to limit:set
 - [`55333ff`](https://github.com/deis/deis/commit/55333ffa13030a54eafda8c61c39ac4e1d4ecb94) README: typo - `ram`
 - [`1950a27`](https://github.com/deis/deis/commit/1950a27b487576676934efaffe8ec7481c8a690a) tests: update README.md with integration test setup
 - [`ed608be`](https://github.com/deis/deis/commit/ed608be80a2a5dae97d5e551e2d57d2f548429aa) (all): make OpDemand the maintainer
 - [`f48b5da`](https://github.com/deis/deis/commit/f48b5dab22c13c9e72ec32dca86c183e48b7ee1c) (all): add load balancer info to EC2 and Rackspace
 - [`d3e7498`](https://github.com/deis/deis/commit/d3e74981c7bafdf169343fe2c8a6e559003641fa) contrib/rackspace: Rackspace details
 - [`ef8bbda`](https://github.com/deis/deis/commit/ef8bbda942299f622bc1fcac1916e24acbb41096) (all): explicitly specify --hosts parameter for clusters:create
 - [`e54b594`](https://github.com/deis/deis/commit/e54b59425f076ca0e1cb546ee12b977d5fe75393) contributing: change "refactor" to "ref"
 - [`1c52eb5`](https://github.com/deis/deis/commit/1c52eb55c9c0b5df26e45cb8e0c9f154e176e74f) contributing: lower commit message length
 - [`f9d7559`](https://github.com/deis/deis/commit/f9d7559ed2447e3feb6c8449afa9e37d1898aacb) contributing.md: import more from docs
 - [`821efa7`](https://github.com/deis/deis/commit/821efa7b72e7b40d414582a23e7e74dc9f762fb3) contributing.md: change styleguide title
 - [`f89b54c`](https://github.com/deis/deis/commit/f89b54cdb1b7b36979f7e8eebb6936ffdeecea4b) reference: add perms to client autodocs
 - [`33c36f2`](https://github.com/deis/deis/commit/33c36f2e431a55ad3bbe881f8dac23d702a30256) client: add more clarity to config:set
 - [`1c42e95`](https://github.com/deis/deis/commit/1c42e95ff430b975306ad550c99f715cb4a4e3ad) client: builds:create no longer coming soon
 - [`d0dbfa9`](https://github.com/deis/deis/commit/d0dbfa9608ffa589cea04d1eb63573d2335af84e) client: the big docstring refactor
 - [`af14a2c`](https://github.com/deis/deis/commit/af14a2cad4a28a90b6f2c51822ca505c7b098434) reference: add config:pull to autodocs
 - [`7c367c1`](https://github.com/deis/deis/commit/7c367c1a3e62a07f4041fe0d85d84fc2d1bad402) (all): Link to ELB timeout configuration
 - [`fc92fef`](https://github.com/deis/deis/commit/fc92fef1c3bf77d1794af86fa7365f352e5c11ad) (all): remove Rackspace support"
 - [`e87bae1`](https://github.com/deis/deis/commit/e87bae1d027a1c31cbc6bbdb56c68c55f4e72781) reference: use deis <command> as title
 - [`a6a3082`](https://github.com/deis/deis/commit/a6a3082c54c1d6278d3b12b467314c7aae0b663b) reference: update reference guide
 - [`5615af7`](https://github.com/deis/deis/commit/5615af7a4228ffb2c9d9d8d4dc4063adf7c7fbe7) reference: update client reference link

#### Maintenance

 - [`c3cbdd1`](https://github.com/deis/deis/commit/c3cbdd1061a6da1048c796f8c513562f6b03460e) CHANGELOG.md: update for v0.11.0
 - [`6941280`](https://github.com/deis/deis/commit/69412805a9eeb7b626fe8c56663532e5c1503344) controller: fix south default values with a data migration
 - [`4f228e5`](https://github.com/deis/deis/commit/4f228e59e7612918de93ee311e1522bb55c4c0b3) controller: update djangorestframework to 2.3.14
 - [`ade588f`](https://github.com/deis/deis/commit/ade588fc1acb7a27c37ec0f7db6851fbefc2eebb) docs: update Sphinx reqs and remove pexpect
 - [`9ef1392`](https://github.com/deis/deis/commit/9ef1392d5318fdaa64f39a867ba38f58198f94a4) tests: update etcd to match CoreOS 402.2.0
 - [`0e6743b`](https://github.com/deis/deis/commit/0e6743b589fea53a629b75fa25f7462072761bfa) controller: update celery to 3.1.13
 - [`c35c92d`](https://github.com/deis/deis/commit/c35c92dcf7418d9c767e8bdfa5f7ced98f60d93f) router: update nginx to stable 1.6.1
 - [`c4a9fa4`](https://github.com/deis/deis/commit/c4a9fa463d0f2a11dc7885cb9051e273ac3c77fe) logger: update go to 1.3.1
 - [`bcffa7c`](https://github.com/deis/deis/commit/bcffa7c9b5445a1e6d6d802524f2c2a26a201c1b) contrib/gce: make GCE script pass flake8
 - [`8572c27`](https://github.com/deis/deis/commit/8572c271ecc0995a809485eec29a9ab87b07ac2c) builder: bump cedarish image
 - [`85419f9`](https://github.com/deis/deis/commit/85419f9a97e45a411382ab70f4d4e07dc2a15769) client: update requests to 2.3.0
 - [`924f0c0`](https://github.com/deis/deis/commit/924f0c0086fb315fe16c1b6bc147f19fc5f17082) client: update docopt to 0.6.2
 - [`4a97424`](https://github.com/deis/deis/commit/4a9742406a756676edf5d4ee9c13574431cc7d2f) controller: update PyYAML to 3.11
 - [`5aca215`](https://github.com/deis/deis/commit/5aca2151f940f560b8991afef99684789f46680f) (all): update CoreOS to 402.2.0; Docker to 1.1.2
 - [`fbfca44`](https://github.com/deis/deis/commit/fbfca44bea9757bedc8237811430c6c94ee811fc) client: remove unused client tests
 - [`796af9a`](https://github.com/deis/deis/commit/796af9ad5f68940f7caae497ccd825c87bf72bab) registry: bump commit
 - [`d0a996e`](https://github.com/deis/deis/commit/d0a996eb35506396f922f2aef75812606df82e55) (all): bump CoreOS to 386.1.0
 - [`44eb938`](https://github.com/deis/deis/commit/44eb93818149def7004ee0c0de754f7b10057458) contrib/coreos: clean up user-data
 - [`eea46ed`](https://github.com/deis/deis/commit/eea46ed794ea7eb754473601ed108edf0cff9074) contrib/coreos: remove unused files
 - [`c2f15a2`](https://github.com/deis/deis/commit/c2f15a20995bc408fea1705d311bcc18a56b26d4) contrib/openstack: clean up OpenStack contrib
 - [`87ed5b4`](https://github.com/deis/deis/commit/87ed5b4801050949c23f01abc88e85acac1e4ae4) router: Use unit file templates
 - [`c629d0c`](https://github.com/deis/deis/commit/c629d0c11218a27dc766ffa3edac08e77d9ceb8b) release: switch master to v0.11.0
 - [`41cd543`](https://github.com/deis/deis/commit/41cd543d09a4da46ed2994ff21ec680048ca2ee0) docs: update CLI versions and download links
 - [`b20ccb2`](https://github.com/deis/deis/commit/b20ccb2d9ef2c44f493ca7a7e9596e96fdce4ec7) contrib/rackspace: bump to CoreOS 379.3

### v0.9.1 -> v0.10.0

#### Features

 - [`46a72ef`](https://github.com/deis/deis/commit/46a72ef7cb5b11d4cca527f49d8ff26fed6220ab) controller: set app dir in etcd
 - [`e376f23`](https://github.com/deis/deis/commit/e376f23e991f207a11c6940a68d263960607df49) builder: add repo check script
 - [`23e1f61`](https://github.com/deis/deis/commit/23e1f617d662772861bc876c8a93ea98e1dec266) builder: add deis build support
 - [`6658dfb`](https://github.com/deis/deis/commit/6658dfb5044849a5f38a159aff7b75bc70e00675) controller: add deis build support
 - [`97da47b`](https://github.com/deis/deis/commit/97da47bea1b836443d775e146521144b40eec6c1) client: add deis build support
 - [`07383e6`](https://github.com/deis/deis/commit/07383e6061ad621060bdbfcd4e414f3417f14803) tests: only pull test-postgresql if it is missing
 - [`cd1c542`](https://github.com/deis/deis/commit/cd1c54248517b2deb28617ef5c835b109db10ab0) tests: only pull deis/test-etcd if it is missing
 - [`57078a4`](https://github.com/deis/deis/commit/57078a47a789d007a47974674c05cccb2c25aa83) controller/coreos.py: Adding TimeoutStartSec to log and announce services
 - [`89cf376`](https://github.com/deis/deis/commit/89cf37640a3acac330d533f12b64bec55e708cdb) (all): move data containers into new unit files
 - [`9ee3e8b`](https://github.com/deis/deis/commit/9ee3e8bdc7f8d777c1184d0ee8dfb4421ce5e79c) logger: add make test/coverage
 - [`34c8997`](https://github.com/deis/deis/commit/34c8997f661e8ce0d8aa83e049b178bd906272fc) (all): allow custom component images
 - [`41ddbab`](https://github.com/deis/deis/commit/41ddbab672f64240194c909b4397b15566a2d968) changelog script: Have script produce markdown, add settings
 - [`4e35d8f`](https://github.com/deis/deis/commit/4e35d8f27fa6d130c577cde47e0ef23adf39d086) contrib/ec2: launch into vpc
 - [`bb50ad0`](https://github.com/deis/deis/commit/bb50ad0e3732470adc5d60bb11654f2c8fc7d07b) contrib: add DigitalOcean provisioner booting CoreOS via kexec
 - [`a0a6b38`](https://github.com/deis/deis/commit/a0a6b38137a711efb302025e9f39fd4b51d876f3) Makefile: move rsync from Vagrantfile to Makefile
 - [`5fbea1a`](https://github.com/deis/deis/commit/5fbea1a7de8f1b476a72d431837d925dbc357030) contrib: Spin up an ELB
 - [`7199d88`](https://github.com/deis/deis/commit/7199d885a4917a46d0202e39b70e8c0bffb9fb8d) docker registry: add Vagrantfile for Docker registry

#### Fixes

 - [`c47ea5d`](https://github.com/deis/deis/commit/c47ea5dea42c7838424165448b8dd7c6ff7c37ef) client: print default string
 - [`b79ba10`](https://github.com/deis/deis/commit/b79ba101403b80e846e1c80b2ae09ff8bbe1ab41) logger: add back message tag
 - [`0fa178a`](https://github.com/deis/deis/commit/0fa178acc5593a777413142b6e60890963a3629f) registry: use confd templated config
 - [`dc55257`](https://github.com/deis/deis/commit/dc552573601c93538c9a8763e8a24285417f21c7) confd: simplify check_cmd and fix grep usage
 - [`5a54a59`](https://github.com/deis/deis/commit/5a54a59e57f433241a808975ac5bc4af21c0adf9) controller: create directory only on create
 - [`31bfe60`](https://github.com/deis/deis/commit/31bfe603d40815f6810646a36034f12ab0eaafa3) controller: use etcd_safe_mkdir
 - [`6ade680`](https://github.com/deis/deis/commit/6ade680a2fdfb115713b17198a88b830f42987bf) controller: create /deis/services on boot
 - [`7f3a5c9`](https://github.com/deis/deis/commit/7f3a5c9cf5bc9d00f8cfd5f308ad71a84688d5db) Makefile: remove bashism in the main Makefile
 - [`2b3c15e`](https://github.com/deis/deis/commit/2b3c15e4e6c779717f3b0958c51974d2f9a57ac9) scheduler/coreos: bump wait_for_announcer timeout to 20minutes
 - [`2d7018e`](https://github.com/deis/deis/commit/2d7018e241aa1499dcfa236d232a63f104bb7366) registry: checkout from a known sha
 - [`c7a22ed`](https://github.com/deis/deis/commit/c7a22ed087c84ade2753be9cf813a056955b7f5e) client: add newline back to `deis scale`
 - [`d91a01a`](https://github.com/deis/deis/commit/d91a01a5637674ebdb44064b0de94c2ca6076bd6) client: initialize app_name
 - [`03279a3`](https://github.com/deis/deis/commit/03279a3ea715ddb61c74bdf0a5d57e805ebd3bcc) controller: install dev_requirements
 - [`182a990`](https://github.com/deis/deis/commit/182a99071acc286635ba1b30e24728c257de8c75) controller: remove tag from source_image
 - [`bf6c249`](https://github.com/deis/deis/commit/bf6c24994e46ea407c1f24e011ad97cd19382ae1) controller: more image parsing fixes"
 - [`7d06fbd`](https://github.com/deis/deis/commit/7d06fbddb94ca2f9f1c7662ad275768d87e2883f) controller: more image parsing fixes
 - [`c415b27`](https://github.com/deis/deis/commit/c415b27ed9b4142ccee212c0971ed7c806c25e6e) controller: build image from source_tag
 - [`01c4b0b`](https://github.com/deis/deis/commit/01c4b0b628e8244bb07d556d6d83345045305e72) controller: correct image name parsing
 - [`2ff233b`](https://github.com/deis/deis/commit/2ff233b3c59a29bd0925467db0017564251dd2d2) builder: revert target_image
 - [`71047e5`](https://github.com/deis/deis/commit/71047e5f668c3fed68799e88ebe77dd3b5f70135) controller: clean up publish_release params
 - [`c6bb251`](https://github.com/deis/deis/commit/c6bb251c232cc15d1bd0e296569f3d07f523d277) controller: add back REGISTRY_URL
 - [`f1f5c8b`](https://github.com/deis/deis/commit/f1f5c8b1d995a45c479bc583a8f0e54fbc308968) controller: import image from remote registry
 - [`e9d1487`](https://github.com/deis/deis/commit/e9d148729b98d896f8178e05b9b6a894ecd9b7ea) registry: fix upstream changes
 - [`84362c9`](https://github.com/deis/deis/commit/84362c911048f2d6479fdbfe722b561aad4e976a) controller: update mock registry
 - [`bf7b223`](https://github.com/deis/deis/commit/bf7b223e9e49911822d11a1b548b63e5a9c360e3) controller: update release process
 - [`92cb26a`](https://github.com/deis/deis/commit/92cb26a054afff117e4c56f434cc00093d23a690) controller: use /bin/sh as entrypoint for run
 - [`3d4bc55`](https://github.com/deis/deis/commit/3d4bc55a03375f049c695d84d27d37124bef42af) client: use present tense when rolling back
 - [`3dbf6d2`](https://github.com/deis/deis/commit/3dbf6d240317c5de9a7a5a7fbfdffdb66d858a46) client: add "but first, coffee" string back
 - [`76c23df`](https://github.com/deis/deis/commit/76c23df5c52ef054afdcb55b0985bddd6673d6bc) controller: use private module
 - [`840f2fc`](https://github.com/deis/deis/commit/840f2fce15fa84157cd33da3d8762496fabbad1e) controller: remove docker engine
 - [`485fb39`](https://github.com/deis/deis/commit/485fb39c16ec7c6bfbb78be7a123ac432c923120) router: update send_timeout from rebase error
 - [`3182ac8`](https://github.com/deis/deis/commit/3182ac870fac07bd3ef9f680793d3c218aaf6c71) router: bump read timeout to 10min for docker build
 - [`e951063`](https://github.com/deis/deis/commit/e9510631c84e512d8504a8a5ab6fc177de68173a) tests: remove unneeded logging from ClearTestSession(
 - [`374ff64`](https://github.com/deis/deis/commit/374ff64ea09e825df81c5fa2d8b6b2f8032ec6bc) systemd units: Fix docker containers being orphaned
 - [`94d59e2`](https://github.com/deis/deis/commit/94d59e2251fcd99b54726b3b46bc1b14b27e11e7) Makefile: ignore status of data services in `make start`
 - [`351027f`](https://github.com/deis/deis/commit/351027f0f2bc708b8f11278f04607691dd8a9492) registry/tests: recreate venv and clean up temp dir
 - [`2507e91`](https://github.com/deis/deis/commit/2507e910d7497f058107cbacee35fe3fd2c26d81) tests: build docker image as first test step
 - [`d3b0783`](https://github.com/deis/deis/commit/d3b07838575d5ea8716f97c6b623896f0a08ce0b) tests: close socket listener in GetRandomPort(
 - [`041e852`](https://github.com/deis/deis/commit/041e8524e7d9c4f7aae9b05ff5ce3ce972e99d3a) controller/registry: add old config vars
 - [`97eddbc`](https://github.com/deis/deis/commit/97eddbc3ba58f02f93a34f7be6761f87565ca3f0) controller: ignore KeyError when purging user from etcd
 - [`0472135`](https://github.com/deis/deis/commit/047213517302c336e6aea81680fbb51d6fe790e9) confd: services check for valid config before reloading
 - [`2e48274`](https://github.com/deis/deis/commit/2e48274ac1644480e29aaea9a8dede756925212a) tests: increase component test timeouts to 20min
 - [`2a4ba73`](https://github.com/deis/deis/commit/2a4ba7306589d12c3f45f36f029f84101ddbe0ca) controller/procfiie: Removed Procfile as its not used anywhere
 - [`2045366`](https://github.com/deis/deis/commit/2045366bf713c8c92ddb27d7d85011b3da502579) contrib/ec2: add script to update the CloudFormation
 - [`94706e7`](https://github.com/deis/deis/commit/94706e7a494c92db5473b0a9fe8df2510a2b07fa) contrib/ec2: add json extension to cloudformation template
 - [`44adc6b`](https://github.com/deis/deis/commit/44adc6b33e7e3fa0237b08f36d57cbdf6b2ef06b) coreos: remove refs to obsolete update-engine-reboot-manager
 - [`9eb4756`](https://github.com/deis/deis/commit/9eb475636e1ba95cb6290924901bdd6110e3ed96) Dockerfile: ensure `apt-get install` is prefixed by `update`
 - [`18c17bb`](https://github.com/deis/deis/commit/18c17bb7927c1b6a463043717c0e4ea806744003) tests: report errors from `docker pull` on main goroutine
 - [`46bca09`](https://github.com/deis/deis/commit/46bca098bce588fb969aaa96ff911edfe8d136ba) Makefile: only install data containers if they don't exist
 - [`5d7c5fb`](https://github.com/deis/deis/commit/5d7c5fbaac0baa53ff566af66cb9babf2bf2c39f) tests: use `docker run --rm` instead of explicitly removing containers
 - [`8fea284`](https://github.com/deis/deis/commit/8fea284c8ee7ddc3f30efb7186191dfec78ecae3) (all): schedule data containers to specific machines
 - [`d5c0d40`](https://github.com/deis/deis/commit/d5c0d40997fb08470c56b53f8d61d3adfdd77054) controller: format command before creating container
 - [`b14228b`](https://github.com/deis/deis/commit/b14228b5b5ff103bd1aca16e303a08c6599a8de2) scheduler/coreos: skip announcer containers for start, stop, destroy
 - [`b4a0865`](https://github.com/deis/deis/commit/b4a086515ed721d896fcb4291068b738982b939d) travis-ci: update logger makefile target to "test-unit"
 - [`0da0f80`](https://github.com/deis/deis/commit/0da0f80e7403780f7b8388797fe9f9acf18ae660) scheduler/coreos: have announcer fail if it cannot get a port
 - [`bda1ed0`](https://github.com/deis/deis/commit/bda1ed098cc7b09343c502cc44d67237892c92e8) scheduler/coreos: correct announce container logic for process types
 - [`f773cde`](https://github.com/deis/deis/commit/f773cde9781561f30d5d6f147f69d868f141c5d3) controller/scheduler: only announce 'web' and 'cmd' processes
 - [`762d15f`](https://github.com/deis/deis/commit/762d15f071180494e0b72b9ce9496e6a64ed1a47) travis: install go cover
 - [`a770a2c`](https://github.com/deis/deis/commit/a770a2c61cc2119ce6f2896e300e78e21cd0eef1) logger: remove Source and Tag from output
 - [`0220c55`](https://github.com/deis/deis/commit/0220c55b40cd18b7d51df5343bc450970b6ecbf3) client: use sys.stdout.write for logs
 - [`daf7600`](https://github.com/deis/deis/commit/daf760043140f833e26301b665341f4b13d2022a) controller: update Dockerfile maintainer
 - [`0ce0864`](https://github.com/deis/deis/commit/0ce0864bc9869e1ee004af6a50108ecfd6bfa4e9) (all): fix etcd_safe_set to not override values
 - [`9209f00`](https://github.com/deis/deis/commit/9209f00632ee0210b7b5e048fec234c70678f49c) contrib/ec2: fix launching into VPC
 - [`68f8a4a`](https://github.com/deis/deis/commit/68f8a4a8f0b429e36ffca2d542ea9cc06a0bd475) contrib: Fix markdown output of changelog script
 - [`6b452e5`](https://github.com/deis/deis/commit/6b452e5f05b56ed6aeb5c9cd64bec49537a8bbd7) services: prevent units from stopping as "failed/failed"
 - [`8d64285`](https://github.com/deis/deis/commit/8d642852875378306a640149c5cf3d1936068e86) Makefile: fix component rsync
 - [`8f3a6c7`](https://github.com/deis/deis/commit/8f3a6c7c88fc0d30c56b125c6faba0af18aea99a) tests: randomize test-database TCP port
 - [`0c59156`](https://github.com/deis/deis/commit/0c59156e49c91aebe144cad5fa7a6cf9b968ce84) tests: use "devicemapper" for STORAGE_DRIVER, not aufs
 - [`44e76f4`](https://github.com/deis/deis/commit/44e76f46fb4073300cdfec45425299560ff97fa2) tests: randomize test port for each service
 - [`a6d7a74`](https://github.com/deis/deis/commit/a6d7a74760d0addd201fa5055a6db875ef5f4f2d) tests: try to scale "cmd" if "web" process type fails
 - [`79cc261`](https://github.com/deis/deis/commit/79cc261b9821dfc70cdc1dfe601ec7e0df9514f2) logger/Makefile: update test-functional target for vendored $GOPATH
 - [`d6e1f2d`](https://github.com/deis/deis/commit/d6e1f2dc87b726ddf71770370b2a03cbbf8d3ff3) registry/test: use env vars for docker daemon connection info
 - [`7865251`](https://github.com/deis/deis/commit/7865251eb57423831554ee3524330dda77e17d92) client: add "Referer" header to requests
 - [`69eb0af`](https://github.com/deis/deis/commit/69eb0af87a765a8bcdd377a90c9da75c4db56f22) docker-registry: up timeout
 - [`93e5797`](https://github.com/deis/deis/commit/93e57979009ef2b3b203aa048712eab2bf714a55) builder: allow filesystems other than btrfs
 - [`4779fee`](https://github.com/deis/deis/commit/4779fee5a7884d2b55651c6b5a0e952887fb3602) Makefile: exit until loop if service becomes failed
 - [`6879f8e`](https://github.com/deis/deis/commit/6879f8e548f3ebf69c5727abbecdc26b3d4d9725) builder: set docker version explicitly
 - [`b93df9a`](https://github.com/deis/deis/commit/b93df9a2f76a4a0aa5ceba47570497595177ac27) scheduler: remove arping hack
 - [`cd9ebaf`](https://github.com/deis/deis/commit/cd9ebafea05c15da0f4343dd17f26cdb9359ce9d) logger: turn off verbose untar'ing
 - [`a3ed364`](https://github.com/deis/deis/commit/a3ed364aac7b43ea4ac3ecf59b588065639bf0f0) controller: disable django admin URL
 - [`eb68a1d`](https://github.com/deis/deis/commit/eb68a1da22d086faf2c5817a18c5c7480d998d8b) controller: disable web UI
 - [`76f2d0b`](https://github.com/deis/deis/commit/76f2d0b5383721369da490455b018359024dc8dc) readthedocs: install CLI requirements explicitly
 - [`0bf1fbe`](https://github.com/deis/deis/commit/0bf1fbe6890c2569d77d03ad99e02578cfa47c8f) vagrant: disable vagrant cachier if present
 - [`db9dc32`](https://github.com/deis/deis/commit/db9dc326b54c660dd347e5f1eeaa6b5227691e21) contrib: check that coreos/user-data has a discovery_url set
 - [`69b077d`](https://github.com/deis/deis/commit/69b077dbe7f107df49406c29d46baf565830d7c3) client: clear cookies on login, logout, register
 - [`fd05885`](https://github.com/deis/deis/commit/fd058859cdcd1bbdf1740dbc3d3ec34fa21058ed) controller: set default release image

#### Documentation

 - [`f15b923`](https://github.com/deis/deis/commit/f15b92382e0dafc741967fdef8ca7e5d6af15e0c) managing_deis: builder needs /deis/services
 - [`74d334c`](https://github.com/deis/deis/commit/74d334c81e4deb5200c76810f352a697d6cc9c0a) using_deis: remove deis builds reference
 - [`a9b13c9`](https://github.com/deis/deis/commit/a9b13c98a0aa6812dbdbbf6687d88ad5ff003f9f) client: elaborate upon builds:create
 - [`f57bf3f`](https://github.com/deis/deis/commit/f57bf3fb911771eb6bc82ffbb8d2cc3e63103970) using_deis: add docs for `deis build`
 - [`e897c8b`](https://github.com/deis/deis/commit/e897c8b7e7ce4cab15787c20338ab1ab72cff2e6) (all): remove Rackspace support
 - [`48b7035`](https://github.com/deis/deis/commit/48b703589c8293cf964d3149e78ab0882eddff4e) (all): add DigitalOcean in a few places it was left out
 - [`a620cf0`](https://github.com/deis/deis/commit/a620cf0f5f391074dd5803fbd58bc0306d5dab1d) ssl: add links for EC2/Rackspace ELBs
 - [`93118d1`](https://github.com/deis/deis/commit/93118d14618d4a709416f5da495726252383d930) installing: add section on installing SSL
 - [`cdf3d6e`](https://github.com/deis/deis/commit/cdf3d6e6c6f952f6128037394be93d90aedcd094) installing: fix up ReST syntax
 - [`46c41eb`](https://github.com/deis/deis/commit/46c41eb73e312c9e9b831d491f050a9df081069f) (all): add back references to DigitalOcean support
 - [`611d199`](https://github.com/deis/deis/commit/611d199b6c99605d10c60d31a0e9b57a453db344) README: fix ram to RAM
 - [`e52d9d7`](https://github.com/deis/deis/commit/e52d9d726c55337379e6b09da0bc211aa9fe1ac4) managing_deis: add dependencies for each component
 - [`f2061f5`](https://github.com/deis/deis/commit/f2061f5d9c4a7da0082881df258430e8be6ae6eb) using_deis/using-buildpacks: fixing unexisting opdemand/example-ruby-sinatra git repository
 - [`d791866`](https://github.com/deis/deis/commit/d79186680204a1907e154b8aed6586f65356e1db) configure-dns: move xip.io to new section
 - [`8a87776`](https://github.com/deis/deis/commit/8a87776a181495b114eee1ccfb2e1c5c73a4954b) installing_deis: add xip.io address for registering on EC2
 - [`23956a1`](https://github.com/deis/deis/commit/23956a165468e4b0a21b725a009b472ea94cc3d7) README: denote `make build` as optional
 - [`6275f79`](https://github.com/deis/deis/commit/6275f79afb3788da782fdcc24bcb7ca1f2413261) contrib: fix DEIS_HOSTS example
 - [`f686c5b`](https://github.com/deis/deis/commit/f686c5b8fc63ca856db42194dde7a9a753baf2c8) (all): fix toctree
 - [`887d30a`](https://github.com/deis/deis/commit/887d30a43f2a0bc18a084cf6e08e49b8824915df) (all): refactor and reorganize
 - [`de06ae4`](https://github.com/deis/deis/commit/de06ae4b225eeea7a87e48e439d03284fb91d49c) pip install: add --upgrade to `pip install deis`
 - [`9fa6a5b`](https://github.com/deis/deis/commit/9fa6a5be546e8fe49f5911b2d44db8cb6715465c) CONTRIBUTING: add allowed types

#### Maintenance

 - [`64708ab`](https://github.com/deis/deis/commit/64708abf5a885cc4ccd2d3de9c49f1c2a44252cf) docs: update CLI versions and download links
 - [`2247a8a`](https://github.com/deis/deis/commit/2247a8a4ec80183a553a7e5ac55b9824d3d938bc) CHANGELOG.md: update for v0.10.0
 - [`6e05dd7`](https://github.com/deis/deis/commit/6e05dd7a748818a6d761d78ea74cc839688f8e9b) registry: bump version
 - [`5814719`](https://github.com/deis/deis/commit/5814719853e02e4d75aa3098e9758ea4fd472850) tests/integration: vendor json library and adjust indentation
 - [`1276365`](https://github.com/deis/deis/commit/12763657367b9af08d682f73206709d736cfb0cd) registry: bump version
 - [`3d4d9f6`](https://github.com/deis/deis/commit/3d4d9f608fbbbdce86c97f74676abcbfc63fa5f9) registry: bump registry version
 - [`03384a7`](https://github.com/deis/deis/commit/03384a77f3aa866eb5ea49deb70d2f440e49fc41) registry: bump version to repository-import
 - [`0d27be5`](https://github.com/deis/deis/commit/0d27be5371d1c72109489df39adf558239d36ba8) tests: update etcd versions in test-etcd Dockerfile
 - [`6acc718`](https://github.com/deis/deis/commit/6acc7181e78b94ede389e687ce8b8d17c862622f) coreos: upgrade to 379.3.0
 - [`5555648`](https://github.com/deis/deis/commit/55556489ffc6fccaeb8d95e9bc91edb50817ff86) contrib/digitalocean: make NYC2 the default region
 - [`f283961`](https://github.com/deis/deis/commit/f28396196e01122b7fb9bb0f36a901e6db3cba63) controller: update South to 1.0 final release
 - [`a5363ef`](https://github.com/deis/deis/commit/a5363ef1e7f6f07d059014252ed5a5a0fd44e3a6) controller: update Django to 1.6.5 security release
 - [`14404e3`](https://github.com/deis/deis/commit/14404e3de67681e881cbfffcd802875aaf911572) tests: update vendored Docker cli code to v1.1.1
 - [`18f05cb`](https://github.com/deis/deis/commit/18f05cb4e973f24fce8fc75977fd76977744c9ff) builder: rely on deis/base to manage confd
 - [`2b1fb73`](https://github.com/deis/deis/commit/2b1fb73363902c1e6617e8a451e6c2b42c8f7f9d) Makefile: clean up rsync excludes
 - [`2b73978`](https://github.com/deis/deis/commit/2b739781fcdf0dfb9629fe7bc5ab8502af00a384) contrib: update Rackspace to 349
 - [`f9ced08`](https://github.com/deis/deis/commit/f9ced08df3ab8e1e764d1d586be132d37b99fb5f) builder: bump to docker v1.0
 - [`1d19f2b`](https://github.com/deis/deis/commit/1d19f2b460a553d65579639c7987567fd559db79) contrib: bump coreos to v349.0.0
 - [`4501c35`](https://github.com/deis/deis/commit/4501c357ca8c5bb23b04f43399fc6ee10a0ddaf3) contrib: bump coreos to v343.0.0
 - [`595a9ef`](https://github.com/deis/deis/commit/595a9ef348b9fdcdaeba4458ebdbcb911306aff3) (all): rely on etcdctl in deis/base

### v0.9.0 -> v0.9.1

#### Fixes

 - [`8d10bf6`](https://github.com/deis/deis/commit/8d10bf6869226d6c9b2ff940358cf4eb773bd225) router: increase timeouts and survive controller/builder failures
 - [`deeec04`](https://github.com/deis/deis/commit/deeec04f5b7a8775a38bb8ce0516cc0a24513040) builder: decouple from controller
 - [`a930b0d`](https://github.com/deis/deis/commit/a930b0db975394d98ed4c38d449ac90481dc825b) router: Add configuration Ingress8000 in deis template for EC2
 - [`d57cc14`](https://github.com/deis/deis/commit/d57cc14b53b1662445f27f954b194865742c77e3) router: Fixes pushing to every router
 - [`7a88613`](https://github.com/deis/deis/commit/7a88613d6d96b6d48ace3ce7053013334e8b402a) client: adjust indentation in keys_add function
 - [`745b26e`](https://github.com/deis/deis/commit/745b26e8de3bd07e72acdd072df0ebecfbb5fff7) registry: fix S3 config options

#### Maintenance

 - [`38c2fb2`](https://github.com/deis/deis/commit/38c2fb2fa057c87ee221ca05f6bb3c637eed2385) versioning: specify image tag in `docker run` and `docker history`
 - [`cef943a`](https://github.com/deis/deis/commit/cef943abfbd99955ad5e801904dbd5fcde19a24a) versioning: update for v0.9.1 release
 - [`79a10fb`](https://github.com/deis/deis/commit/79a10fb61d367fa1660004db4138195f30d64449) CHANGELOG.md: update for v0.9.1 release
 - [`afc27c4`](https://github.com/deis/deis/commit/afc27c4cd667fdb326ca914c05a9ea42a26c6ee5) contrib: Removed unused public port

### v0.8.0 -> v0.9.0

#### Features

 - [`8428bd5`](https://github.com/deis/deis/commit/8428bd5a7c2db6ad6246a400ce5d28dbfc477529) router: proxy builder through router
 - [`bddb43e`](https://github.com/deis/deis/commit/bddb43ee2ba5dadb9c63a1548d50bab265482c40) user-data: adds nsenter alias
 - [`5ad1fc0`](https://github.com/deis/deis/commit/5ad1fc07f240f684ee78774e1d5eb9ed22ff11b7) router: route deis.domain to controller
 - [`5d60551`](https://github.com/deis/deis/commit/5d60551f07fa64ebf93deeca21373ea32ca03f85) vagrant: make rsync as default
 - [`e86d9d8`](https://github.com/deis/deis/commit/e86d9d8a9b046be6dd99479d6029d60cfde25acb) builder: build apps from more than master
 - [`321b96c`](https://github.com/deis/deis/commit/321b96c9d6698ad75b6263d53a7ecf69377ce4d3) controller: allow shared users domain access
 - [`65bf80e`](https://github.com/deis/deis/commit/65bf80e3163142f787e6b5ff3fd75de3a1e0a681) controller: log domains deployed
 - [`c1c7ee7`](https://github.com/deis/deis/commit/c1c7ee76d0b9f72dc2a4d46b087055cacdc31c08) controller: hook up domains to router
 - [`23df4ab`](https://github.com/deis/deis/commit/23df4abca42cf43c796b33536dbe8f724c7b15c7) controller: add Domain model
 - [`fcadd48`](https://github.com/deis/deis/commit/fcadd48c6d179977b11f1414a623a7a7772c0812) controller: Toggle registration using etcd
 - [`730149c`](https://github.com/deis/deis/commit/730149cac7957312356b2005794fe353021256ec) builder: update builder to pass SHA/Procfile/Dockerfile (if available
 - [`c5aec20`](https://github.com/deis/deis/commit/c5aec20be5657af84def14671adbf45dd3ff0963) dockerfile: improve dockerfile and procfile workflow
 - [`f330bfd`](https://github.com/deis/deis/commit/f330bfd2271a31726c1c3ed1b69dd07847e25062) router: support for multiple routers
 - [`f845dc8`](https://github.com/deis/deis/commit/f845dc8fc4a7eeee23b00ab056fadae5b3f065ae) router: configurable gzip settings
 - [`e0ba0a9`](https://github.com/deis/deis/commit/e0ba0a976e0734d7fab069be1a0610d2f8078ddb) builder: expose runtime configuration during slugbuilder execution
 - [`097827c`](https://github.com/deis/deis/commit/097827c46c43d50e571f0a8563c9ab9b63d4c0d9) Vagrantfile: add fallback to rsync
 - [`2bce3f4`](https://github.com/deis/deis/commit/2bce3f4eb5ae6d877bf6f3541c9b8171d12c6aeb) contrib: add changelog script
 - [`4712658`](https://github.com/deis/deis/commit/47126583d95176e5a99eea823eb8b6df98dd0641) contrib: add CoreOS flair to motd
 - [`11b07e2`](https://github.com/deis/deis/commit/11b07e2e6d6bc068ba215c7d506778fb60684da2) contrib/coreos/user-data: create Deis motd
 - [`4808ea9`](https://github.com/deis/deis/commit/4808ea90dfe9722a58986f8ec583e77384bff0fb) client: add --app option to apps:run
 - [`b2e2b78`](https://github.com/deis/deis/commit/b2e2b78e17886d0390fa9f67bf63ba368b726d33) client: add auth command

#### Fixes

 - [`0470a38`](https://github.com/deis/deis/commit/0470a3834961667ffa12527089bd1220017d8f7e) controller: add source release version
 - [`d290bb2`](https://github.com/deis/deis/commit/d290bb2e41594d075afd2ce4426456e32a3a7890) controller: revert "<no value>" polling target to confd_settings.py
 - [`733a61c`](https://github.com/deis/deis/commit/733a61c4a1a208dbf420d682324ae5953e4fda31) docker: always use ":latest" tag in docker pull
 - [`b800e9b`](https://github.com/deis/deis/commit/b800e9b34a29b4b5e9236a14cd6df6f5c554df9b) router: restore check_cmd and grep nginx.conf
 - [`f349741`](https://github.com/deis/deis/commit/f349741525f54ffbeb35baf63f2624ce83f28373) controller: change timeout for docker container template
 - [`2767c47`](https://github.com/deis/deis/commit/2767c47310725b6c83a19c2db59bd46560c886ef) travis: add dummy FLEETCTL_TUNNEL envvar
 - [`c0eae21`](https://github.com/deis/deis/commit/c0eae21d93a2b6e78296b227db8356d748426e2b) Makefile: enable targets to work for non-Vagrant
 - [`07d49ef`](https://github.com/deis/deis/commit/07d49ef782248dc8effaef56c5142b0d382ed141) controller: watch fifth column for state
 - [`82a0339`](https://github.com/deis/deis/commit/82a0339c467ce10f9d5c2f155eda1b69b86054b4) controller: use new `fleetctl` semantics
 - [`1c67726`](https://github.com/deis/deis/commit/1c67726e33156399b31b331d1ff1154935dd38d7) data-containers: change test command to "inspect"
 - [`52b93ac`](https://github.com/deis/deis/commit/52b93ac1f65582cee28e78cad8ef88b7ac1e772b) Makefile: update to work with fleetctl 0.3.2
 - [`91e7e21`](https://github.com/deis/deis/commit/91e7e214a315247b0c3332b6a3cc669d2967d8bc) test: update example-ruby-sinatra location
 - [`77c1703`](https://github.com/deis/deis/commit/77c17034f4efa9fc71ce6621fa316073ecddbe72) contrib: revert seed hosts for CDN workaround
 - [`3e1708c`](https://github.com/deis/deis/commit/3e1708c19b2ac5e421406663ad2d9f84797586bb) Makefile: Information about missing FLEETCTL_TUNNEL
 - [`f4a41ad`](https://github.com/deis/deis/commit/f4a41ad9176b20d3380120f80393f1a8cf4128fc) Makefile: include router component in `make build` and friends
 - [`1fdbf2c`](https://github.com/deis/deis/commit/1fdbf2c22510c749a7e6fd97dd2db3f0d23a6692) Makefile: use fleetctl to filter actual deis-* units
 - [`129ad5d`](https://github.com/deis/deis/commit/129ad5d2da771a3cc74a6f4c7dcfc41171694f43) controller: change api logs response to 204 when no logs
 - [`bdf08ba`](https://github.com/deis/deis/commit/bdf08ba57fe3ed17cad7686c3d8c1fc2bacfc64b) signals: use different dispatch_uid for log signals
 - [`6aad59c`](https://github.com/deis/deis/commit/6aad59c266fdae9e0ce4722e1d5a93616e51d149) (all): set etcd keys safely in /bin/boot
 - [`968efd0`](https://github.com/deis/deis/commit/968efd0db4ae788a5f3017399b359d75f4c453a0) router: create /deis/domains
 - [`6886795`](https://github.com/deis/deis/commit/688679565119c054d9fe65fba1a0f4b2780b4f51) controller: give access to /domains
 - [`dbfed01`](https://github.com/deis/deis/commit/dbfed019296b48f4bdf531e8866ab8e74e86b70c) client: handle created response code on domains:add
 - [`9a9c40c`](https://github.com/deis/deis/commit/9a9c40cce25ff54708dcbbe1df51d27a4de36aef) controller: do not deploy releases on domains
 - [`526cd06`](https://github.com/deis/deis/commit/526cd0658b079a7f6c1d999464678cbe049fc76f) client: cleanup domains output
 - [`ed642ff`](https://github.com/deis/deis/commit/ed642ff707da54a13fdf2c71bcb184514f9f88cf) router: resolve golang template scoping issue
 - [`1646f72`](https://github.com/deis/deis/commit/1646f72d6df25a554ae651be353dd67a6738ea8d) controller: add /domains access in confd
 - [`133f707`](https://github.com/deis/deis/commit/133f7071b12d152a852e2a23e5d6bddf732117ec) controller: revert e647336741
 - [`64c90e7`](https://github.com/deis/deis/commit/64c90e7101243ad6cf691ec5096fd91b061bc73e) confd: give controller/router /deis access
 - [`1c06b51`](https://github.com/deis/deis/commit/1c06b51567e69227796833f976229d6bd3f3ee74) controller: cast registrationEnabled to bool
 - [`2dee62c`](https://github.com/deis/deis/commit/2dee62cc1a9f6579d55884c638ec08f13622d126) scheduler: detect exposed ports from image for PORT envvar
 - [`3d5d243`](https://github.com/deis/deis/commit/3d5d243dfaada45014e91df0058cbed8a8a248f1) tests: close db connections during threaded execution
 - [`5710ccb`](https://github.com/deis/deis/commit/5710ccb730b1b46d56c0bbfadfe424e93ecf218d) flake8: resolve code formatting issues
 - [`f680df6`](https://github.com/deis/deis/commit/f680df698de17f073a72d99b896be711c2155c5f) deploy: rework initial deploy logic based on workflow
 - [`877a87b`](https://github.com/deis/deis/commit/877a87bc6de6e9e5c0c67d3a23be4cbf48942899) builder: include default_process_types from .release
 - [`b4ff75e`](https://github.com/deis/deis/commit/b4ff75e275b83e0f8c5384857b7f4e060acb70cc) migration: squash migrations, add newline
 - [`cdbd9b9`](https://github.com/deis/deis/commit/cdbd9b975bdda52c9118ef2945874e007e518c0b) scale: move initial scaling logic to model, add tests
 - [`553978c`](https://github.com/deis/deis/commit/553978cc3e800cce48dbb9ab92d947c9eca98444) controller: only scale web=1 if procfile has web entry
 - [`2984581`](https://github.com/deis/deis/commit/29845814f94767f7dee69cc924f22effe939da91) schema: remove uniqueness constraint on containers
 - [`c8efffd`](https://github.com/deis/deis/commit/c8efffd0ed4a5fd1b8421cd54871fd67aa886519) announcer: detect first exposed port instead of hardcoding to 5000/tcp
 - [`a056894`](https://github.com/deis/deis/commit/a056894cde01146c395573f0d652e2fe74e2f2ca) test: remove `-K` ssh-add option that isn't available on Linux
 - [`1d46330`](https://github.com/deis/deis/commit/1d463306cdc89fb8e1f4e3c6d36a0cf50577dfa6) builder: add scale(
 - [`fd88251`](https://github.com/deis/deis/commit/fd88251a7ac12ee1612ac5d6d1354ac6f7624e4a) flake8: fix spaces and a docstring
 - [`95c2c77`](https://github.com/deis/deis/commit/95c2c77e0555d2b14b0212529553138bd08a10e1) app: only set initial structure when build hook is triggered
 - [`c951ffa`](https://github.com/deis/deis/commit/c951ffab4c4ce9e82ad28c90555a7d23316bf219) Vagrantfile: fix vagrant sync command
 - [`96e0428`](https://github.com/deis/deis/commit/96e042853e1fec7e595556cd1741cad689a4008d) README: fix typo
 - [`d3dfa37`](https://github.com/deis/deis/commit/d3dfa376d86ffcbbbcab8211edb161b4cfde2d72) contrib: seed hosts for CDN workaround
 - [`166282a`](https://github.com/deis/deis/commit/166282a03bc6157dad4dbdd6e0ff5ffd343f4314) CHANGELOG: prep for changelog script
 - [`90c73d0`](https://github.com/deis/deis/commit/90c73d0bcae1ff9fec5cafc7392347dddbb1bafe) test: don't fail if already registered, and pause for builder
 - [`48cd94f`](https://github.com/deis/deis/commit/48cd94f9259193bf232547d6416b47b7e7616166) client: return 1 on two more error cases
 - [`45e5e0c`](https://github.com/deis/deis/commit/45e5e0c89eb7964e0463a1e0117a01538b088234) contrib: correct coloring
 - [`23f0f0c`](https://github.com/deis/deis/commit/23f0f0c48ee3ed08c7f285547647c13e78fa6eeb) logger: allow digits in app name when parsing log file
 - [`33f72a1`](https://github.com/deis/deis/commit/33f72a175a720663d3b8f0e97c7bd7c6c014cc14) controller: restore App.destroy(
 - [`998498f`](https://github.com/deis/deis/commit/998498f67ea52aeef5af21416767e01c7fb9cdb5) Makefile: fix 'check-fleet' rule for $FLEETCTL_TUNNEL with port.
 - [`6d64cbf`](https://github.com/deis/deis/commit/6d64cbf2ce09f92f45563af48e07b75ac6c13432) readme: unify language with deis.io
 - [`c510fc5`](https://github.com/deis/deis/commit/c510fc5971f09b3a886e507e2867b2db36d704e6) controller: give processes unique counters
 - [`2c699ed`](https://github.com/deis/deis/commit/2c699ed67daf7c539b9b1940d989d8d9909e3704) controller: remove username length validation
 - [`75aeb67`](https://github.com/deis/deis/commit/75aeb670fa7d74870e02b88c28ad271a23582ee3) registry: store apps as appname
 - [`926c8f4`](https://github.com/deis/deis/commit/926c8f49937bb21f2197a1388fc907e68fe34a20) client: return non-zero on errors
 - [`a0cf80c`](https://github.com/deis/deis/commit/a0cf80ce694ec29ce980baffad409662dec73a64) docker: use btrfs inside builder
 - [`15751ad`](https://github.com/deis/deis/commit/15751ad0a45b3be074512d21fa92378520ac8067) builder: do not check for a Procfile
 - [`4edfb2c`](https://github.com/deis/deis/commit/4edfb2ce0a3805f5f6e7d0c4d5f9a7b2fef64b58) pip: update location of pip installer

#### Documentation

 - [`70fe559`](https://github.com/deis/deis/commit/70fe559cb5edd49f23aab7507bc2e263519c0ca0) (all): document DEIS_NUM_ROUTERS
 - [`ad8eabd`](https://github.com/deis/deis/commit/ad8eabd34423fbb19a707cd6a4f1c4f10a2e2b32) readme: remove cruft and add clarifications
 - [`ccdfc34`](https://github.com/deis/deis/commit/ccdfc34f6a5602432f55b94c17ac94224a9e0819) contrib/rackspace: document manual update for CoreOS"
 - [`d759e1e`](https://github.com/deis/deis/commit/d759e1ec348201503a30da5f5ef09579c05fb6cd) README.md: Update Documentation for creating Applicaiton
 - [`e98d9c2`](https://github.com/deis/deis/commit/e98d9c224650a9a58f4ccfc56a06054278ee2ffd) releases: update for Docker Index tagging capability
 - [`d459eec`](https://github.com/deis/deis/commit/d459eec3c8d3fe8e3bc02e6e3ebb2ea7066bdfa5) configure-dns: fix typo
 - [`d4aa8c9`](https://github.com/deis/deis/commit/d4aa8c95715e3e03b439d7ac65f08240793f20de) UPGRADING.md: add upgrade documentation
 - [`20ecc29`](https://github.com/deis/deis/commit/20ecc297e8e9e62a4b3ef9c48d82955123258c12) contrib/rackspace: document manual update for CoreOS
 - [`54888af`](https://github.com/deis/deis/commit/54888af304455ae79a94d668b1f6b216981ab629) developer: use note styling
 - [`dbbb7b4`](https://github.com/deis/deis/commit/dbbb7b4b7e425f8bc5309c720812431b1a9720d0) developer: add domain documentation
 - [`87159d5`](https://github.com/deis/deis/commit/87159d56a5e5f4f6972de296ae73f14337108e5b) developer: flesh out developer docs
 - [`ae0365b`](https://github.com/deis/deis/commit/ae0365b35148959a428ab097b80cff9464b54b93) (all): update docs for local3.deisapp.com, local5.deisapp.com
 - [`db2d845`](https://github.com/deis/deis/commit/db2d845b21b4191f9e73e3da0ebaa668a7407c7c) bare-metal: filled in bare metal section in provision doc
 - [`ff3e029`](https://github.com/deis/deis/commit/ff3e02959bbdf8e1a0e7811963190b87cc832243) (all): add bare-metal guide
 - [`9435bfd`](https://github.com/deis/deis/commit/9435bfd2074850151a7fd1e17720212b55319f08) README: add troubleshooting for 'not reachable' peers
 - [`f3c60cc`](https://github.com/deis/deis/commit/f3c60cca8ab0d9e00b05cf5672b2217707d02b2e) contrib: github-changes -> changelog script
 - [`b70aa56`](https://github.com/deis/deis/commit/b70aa56c00694efef55ca0dad39f5933344c7e9c) (all): clarify clusters:create options
 - [`ac15bce`](https://github.com/deis/deis/commit/ac15bce3c72b6c9defd4056960c93750af2f2c9c) README: adds troubleshooting information
 - [`28f3727`](https://github.com/deis/deis/commit/28f37273b65d47682fdc33cd6d53c94ed590c6fa) (all): strongly suggest Deis clusters of 3+ nodes
 - [`900e583`](https://github.com/deis/deis/commit/900e5831eca8b82f703d9b3ca0bf570b7bdc06ec) contributing: add TESTING footer
 - [`842838a`](https://github.com/deis/deis/commit/842838ae397b0fa1935742f88e4afd29c454c58c) vagrant: update references to vagrant files
 - [`2dd299c`](https://github.com/deis/deis/commit/2dd299cdc548e191c22f0a67bedfa6eb6cf178c7) README: install nfs on Ubuntu
 - [`98b478a`](https://github.com/deis/deis/commit/98b478a36ad490b044db6088e6eb229889f88fc8) contributing: split out DCO

#### Maintenance

 - [`dc7bdda`](https://github.com/deis/deis/commit/dc7bdda1c192bf7695936514e2d4d3c5c5e88f23) (all): update Docker tags to v0.9.0
 - [`3084115`](https://github.com/deis/deis/commit/3084115563e48dab5892c366553e344b9812a7af) client: update deis CLI binaries at S3 to v0.9.0
 - [`072ac1a`](https://github.com/deis/deis/commit/072ac1a533530b1f4fbf2e6bfa4e089535e4b68c) contrib/(all): update to CoreOS 324.2.0
 - [`369c2f4`](https://github.com/deis/deis/commit/369c2f431c3a66a52b35af5299f3c916f153762c) contrib/(all): update to CoreOS 324.1.0
 - [`7ccc031`](https://github.com/deis/deis/commit/7ccc031cf47ab8cc3119c67213461284fd791ef3) controller: remove unused coveragerc
 - [`1333f6a`](https://github.com/deis/deis/commit/1333f6a492956617a4443e4274bdc6c7220725c9) builder: remove unused authorized_keys
 - [`d72013a`](https://github.com/deis/deis/commit/d72013a2204462d831cce77bf2abf887f7c1aa64) controller: update redis to 2.9.1
 - [`8556509`](https://github.com/deis/deis/commit/8556509b5f3c4642dd762cd8c0bd863314d084ab) controller: update celery to 3.1.11
 - [`a088a68`](https://github.com/deis/deis/commit/a088a6805ce481e3889e65b54f640b7f74ab36d0) controller: update djangorestframework to 2.3.13
 - [`5e418ee`](https://github.com/deis/deis/commit/5e418ee014f160b149a01e58d262c3a81dc0f30b) contrib/(all): update to CoreOS 310.1.0
 - [`8347cbe`](https://github.com/deis/deis/commit/8347cbe641168d5061ce01eb53d80d277b00de22) builder: update docker to 0.11.1
 - [`7925589`](https://github.com/deis/deis/commit/792558964bbcb31413a811a742da6832063006cd) Vagrantfile: update to CoreOS 310.1.0
 - [`c0f7b74`](https://github.com/deis/deis/commit/c0f7b74aa4748c9dfa936793ea60917b0edf9530) controller: install new pip 1.5.5
 - [`0f50a11`](https://github.com/deis/deis/commit/0f50a113e9d61970c2c65729c310294e79481067) controller: remove unused gevent package


### v0.8.0 (2014/05/06 18:30 +00:00)
- [54ff9ca](https://github.com/deis/deis/commit/54ff9caaaaa2be20960afe41c883014a796b81b2) Updated CLI binaries and links in client README. (@mboersma)
- [4694844](https://github.com/deis/deis/commit/4694844e603583d2b0cee0fa0b4922a9857dc5d9) Switch master to v0.8.0. (@mboersma)
- [#715](https://github.com/deis/deis/pull/715) Merge pull request #715 from opdemand/filesystem-docs (@opdemand)
- [469ee72](https://github.com/deis/deis/commit/469ee72a8e180485fbe7a394f052b1a47f3ae234) refactor(scheduler): use coreos/fleet for container scheduling (@gabrtv)
- [aa8bafc](https://github.com/deis/deis/commit/aa8bafc5608862f220f7fea245b2a6e27141fa91) refactor(scheduler): adapt web app to new cluster domain object (@gabrtv)
- [b5aea48](https://github.com/deis/deis/commit/b5aea482dd567ba8bfe4c4a9009bfe18e7afeec1) refactor(cm): remove deprecated cm and provider packages from test and coverage (@gabrtv)
- [10c8c4e](https://github.com/deis/deis/commit/10c8c4eaa428bc9a691931e5e315ba57882a824e) refactor(client): changes for scheduler refactoring (@gabrtv)
- [51f5d6d](https://github.com/deis/deis/commit/51f5d6d57bc80c882c07ef3c649892bcbbdfed10) chore(builder): upgrade to docker 0.9.1 (@gabrtv)
- [30f1651](https://github.com/deis/deis/commit/30f1651c8939db20920f87d0b92912cd8b8748cf) feat(docker): deis meta project (@gabrtv)
- [d4e13a4](https://github.com/deis/deis/commit/d4e13a44961010ac9776bc3878d434537b09d53c) refactor(scheduler): coreos scheduler implementation (@gabrtv)
- [2c72ff2](https://github.com/deis/deis/commit/2c72ff23d0c2333fe9359b06bc750cc37c16b561) refactor(builder): improve build + config = release (@gabrtv)
- [e1ccc29](https://github.com/deis/deis/commit/e1ccc29c4c72458be0683dc0442f0a99a50f5740) refactor(builder): remove slug metadata from builder (@gabrtv)
- [a20798a](https://github.com/deis/deis/commit/a20798aa4b91919c2800d9871d0a9e9c594e0576) feat(scheduler): add command handling for procfile-style execution (@gabrtv)
- [0a217fc](https://github.com/deis/deis/commit/0a217fcdd83bac96543a29d60958e974907992e8) refactor(vagrant): new CoreOS Vagrantfile (@gabrtv)
- [8110aaf](https://github.com/deis/deis/commit/8110aaf91b876ea885c96a5031fdc795c4cd5489) refactor(makefile): dispatch to relevant coreos utilities (@gabrtv)
- [4220183](https://github.com/deis/deis/commit/4220183ff633b77d297c0e65363e4fb8bd388a0c) feat(registry): add private docker registry (@gabrtv)
- [8d9f314](https://github.com/deis/deis/commit/8d9f31400562a4abce66c0c9a1bfb0690d89f301) feat(router): add nginx router driven by confd/etcd (@gabrtv)
- [0a081cd](https://github.com/deis/deis/commit/0a081cdc0248ada1f3fe0cb489b085fb471e7f8f) chore(ruby): remove ruby files (@gabrtv)
- [72ddbaa](https://github.com/deis/deis/commit/72ddbaa5a7f6e7a1fa4d5ed3ff54e2f0bf2551a5) chore(docker): remove registry from deis meta project (@gabrtv)
- [bc422ef](https://github.com/deis/deis/commit/bc422ef605efd25b6f7c2863f3c5d2b1ad2e22a0) chore(builder): show `docker push` output from builder for debugging (@gabrtv)
- [98387cf](https://github.com/deis/deis/commit/98387cfb172aab60c14e19f4de65ab72a263e8cb) chore(coverage): add .coveragerc in controller root (@gabrtv)
- [680f1a8](https://github.com/deis/deis/commit/680f1a8b1ff34237d9a0032a12644196c2536f8f) chore(router): expose router on 80 and 443, add timeout for image download (@gabrtv)
- [dd065bd](https://github.com/deis/deis/commit/dd065bd054ffcb373c258bd4f95df75da4880208) refactor(client): make apps_open cluster aware (@gabrtv)
- [039ad59](https://github.com/deis/deis/commit/039ad5916124cd0c04a998253570de8deb97ab70) refactor(deis-run): new `deis run` implementation and celery task (@gabrtv)
- [ac956d1](https://github.com/deis/deis/commit/ac956d1ee28091c211e688126f4a98ab1c65ecaa) refactor(registry): new registry package with mock and private implementations (@gabrtv)
- [ba19ac2](https://github.com/deis/deis/commit/ba19ac2385130eb77e93aa8bca9b3045523fd0ff) fix(boot): dot remove docker.pid, fix sock test (@gabrtv)
- [e07b18d](https://github.com/deis/deis/commit/e07b18d8fb0f13c11ac27d029a73c4b7dae60e39) fix(cache): change cache workdir to /var/lib/redis for read/write filesystem (@gabrtv)
- [426fbce](https://github.com/deis/deis/commit/426fbce728a92c54dc999b6394a408c115042b4b) fix(boot): only remove docker.sock if it exists (@gabrtv)
- [01f2225](https://github.com/deis/deis/commit/01f222551cb1b82594c6fb22de8dec4096cf41f8) refactor(bin/boot) consolidate into a single script with debug and optional publishing (@gabrtv)
- [41645ee](https://github.com/deis/deis/commit/41645ee2463f5645f5569b0beacf84e58c454f65) doc(readme): update README instructions for dev setup (@gabrtv)
- [37964ae](https://github.com/deis/deis/commit/37964aeff8efa521c214fc00dae323a2c8d9a2a5) chore(builder): add debug, move confd loop before daemon execs (@gabrtv)
- [f46834a](https://github.com/deis/deis/commit/f46834ab1e76b9ecf03984906cef3da486df9f39) fix(builder): add bridge CIDR address to fix etcd connectivity (@gabrtv)
- [162928c](https://github.com/deis/deis/commit/162928c11610d276c2e3614294e7a77e51ca7f43) chore(systemd): change units to publish on known ports (@gabrtv)
- [802666b](https://github.com/deis/deis/commit/802666b8fee0ae307780c2fd287d353ca8dca599) chore(scheduler): disable auto router and logger on cluster creation (@gabrtv)
- [be7a003](https://github.com/deis/deis/commit/be7a00361c063fc16a9a18eb805abab92e62f550) chore(docker): move deis meta project to contrib/docker (@gabrtv)
- [c08eaf8](https://github.com/deis/deis/commit/c08eaf839943cafaaeff55f426c2f9f1b253691a) chore(user-data): change yaml content to literal block (@gabrtv)
- [18d1626](https://github.com/deis/deis/commit/18d1626d15117f41c086205528fee72e5a8dc314) fix(vagrant): switch vm to 2gb of memory (@gabrtv)
- [b7f29fa](https://github.com/deis/deis/commit/b7f29fad0a69e52d635a004332b8969327a4a272) chore(router): disable check_cmd for now (@gabrtv)
- [bb67e58](https://github.com/deis/deis/commit/bb67e581eb468b0239409b26397cc2f0ad4481ab) chore(registry): turn logging to info (@gabrtv)
- [6e1d70d](https://github.com/deis/deis/commit/6e1d70d98607bc578df75d917faaf6cf9843fe8d) chore(database): remove conf.d comments (@gabrtv)
- [774ef63](https://github.com/deis/deis/commit/774ef6303b44c0b6d58c2f600bc9fd8391cfe9a2) fix(vagrant): add host ip to /etc/hosts (@gabrtv)
- [0e2b42b](https://github.com/deis/deis/commit/0e2b42bd025f785af5130197f110089257b8896c) chore(makefile): remove docker daemon from logs (@gabrtv)
- [34aab76](https://github.com/deis/deis/commit/34aab76228596db5eb95114abcbbbb5f2b42994d) chore(vagrant): bump memory to 4096 (@gabrtv)
- [6c15217](https://github.com/deis/deis/commit/6c15217c5e13d32c99ccded949504a10d3d91eba) fix(router): handle non-zero rc from etcdctl mkdir (@gabrtv)
- [dfc762c](https://github.com/deis/deis/commit/dfc762ca1923ffc5aa9dd4608730e07f7b899159) fix(scheduler) switch to getent for host resolution (@gabrtv)
- [2b9fa82](https://github.com/deis/deis/commit/2b9fa8292d99556db465d4155f73c99e7ca3fe0c) fix(etcdctl): use no-sync option to prevent connectivity issues (@gabrtv)
- [0b54a13](https://github.com/deis/deis/commit/0b54a135a54bc3b58d9a179a1395a3e82bc57a76) fix(docker): switch to aufs storage backend (@gabrtv)
- [7e9b418](https://github.com/deis/deis/commit/7e9b41853f74c23ca23d533f7258d5e2dee9a4d3) fix(systemd): wait for etcd.service (@gabrtv)
- [6aa0092](https://github.com/deis/deis/commit/6aa0092df404c127eef85af7e52638b902ea2e54) chore(journal): suppress etcdctl output, add other clarifying output (@gabrtv)
- [e173103](https://github.com/deis/deis/commit/e173103aa86759ff6d989192603876b3b65e099a) perf(images): seed docker images as a oneshot (@gabrtv)
- [4029205](https://github.com/deis/deis/commit/4029205c775b6518969aaf14810cfb1c30232324) fix(scheduler): remove unused functions (@bacongobbler)
- [853ed89](https://github.com/deis/deis/commit/853ed8970a441c27dd979ac43bd1b47a93fdb06c) docs(readme): `make pull` of cached layers to speed `make build` (@gabrtv)
- [6efe264](https://github.com/deis/deis/commit/6efe2640228471b85c3542ba384f4294ca11648a) perf(images): seed the deis-registry with slugrunner on boot (@gabrtv)
- [decd712](https://github.com/deis/deis/commit/decd71211b8b900780822e7e4cae73b6b708f446) fix(builder): fix typo in confd toml (@gabrtv)
- [163359a](https://github.com/deis/deis/commit/163359a61d55e99b3c399b332f7c2b6da4aaa1b4) feat(builder): scale the web process by 1 (@bacongobbler)
- [f849c47](https://github.com/deis/deis/commit/f849c4765d5949817f1bfbd8a2ee0270d31cfb07) fix(docker): use btrfs as coreos does. (@mboersma)
- [2fb3dfd](https://github.com/deis/deis/commit/2fb3dfd86ac1433892ed64cec6e5185b0a877721) refactor(controller): remove chef gem (@bacongobbler)
- [9e11148](https://github.com/deis/deis/commit/9e111489a7458d4567715ab5b384be5b64035a31) fix(tests): use TransactionTestCase where threads are used (@gabrtv)
- [17d36e1](https://github.com/deis/deis/commit/17d36e1cd7c2f9eb02576d15d0236e10fdf5731b) fix(controller): use postgresql for tests (@bacongobbler)
- [42e45a6](https://github.com/deis/deis/commit/42e45a6a14c540ad8d083dd46d8aad9ffa82a106) fix(controller): allow only admins cluster access (@bacongobbler)
- [5b3d268](https://github.com/deis/deis/commit/5b3d268062834455d20191e272b297f72585b79f) bug(contrib/coreos): fix Docker bridge (@carmstrong)
- [231ede3](https://github.com/deis/deis/commit/231ede3816c3d0bb8576ca46d7a7d2c453f5c390) fix(contrib): update ec2 cloudformation docs (@bacongobbler)
- [110f093](https://github.com/deis/deis/commit/110f09360213975ec13b577f4b556e01e15333ba) bug(docker): bind docker bridge to a specific IP (@carmstrong)
- [741825c](https://github.com/deis/deis/commit/741825cc3588f79ad83dc4daa8f197e7069d6d26) bug(systemd): fix dependencies (@carmstrong)
- [48c06f7](https://github.com/deis/deis/commit/48c06f700b0e58a44d0106006c93ed8995da4859) fix(docs): remove outdated autodocs (@bacongobbler)
- [fd9b42e](https://github.com/deis/deis/commit/fd9b42eff5e2d3e8f7bbb44c8284da6aa3e3121f) fix(client): add clusters command to `deis help` (@bacongobbler)
- [8952179](https://github.com/deis/deis/commit/8952179cd5544f18af98be1dffcde32ff444b087) fix(*): make containers host-aware (@bacongobbler)
- [75ed9d7](https://github.com/deis/deis/commit/75ed9d71c95b4e3398f6464b98d2c01dd09d0619) docs(scheduler): update documentation for CoreOS/Fleet/cluster changes (@mboersma)
- [cfafc0d](https://github.com/deis/deis/commit/cfafc0dbb921ccba26630ab0826dd6c43b5bbc44) refactor(discovery): remove discovery component (@bacongobbler)
- [57cde57](https://github.com/deis/deis/commit/57cde575651a094c4273df201b47c6579c7ec437) feat(contrib/vagrant): support multiple-VM environment (@carmstrong)
- [4111501](https://github.com/deis/deis/commit/4111501dc9de0c99fdda42d82b822d989d6d382f) refactor(contrib): update Deis deployment on AWS (@bacongobbler)
- [a0f3437](https://github.com/deis/deis/commit/a0f3437cd0f7897b399174b2bb049e879b8ce267) refactor(controller): optimize docker build layers (@bacongobbler)
- [cd6a38e](https://github.com/deis/deis/commit/cd6a38e798fbd1ef60efa518efbee473d8d9e580) feat(controller): add container state via FSM (@bacongobbler)
- [e1eafbc](https://github.com/deis/deis/commit/e1eafbc516ec9c2fc6f237b575f5ccff6d421e32) fix(registry): chown /data for data container (@gabrtv)
- [418c025](https://github.com/deis/deis/commit/418c0253c03ea59f91c441e9a129c9afebf7b7e9) fix(controller): install git (@bacongobbler)
- [06659ec](https://github.com/deis/deis/commit/06659eccc1cc9be645d4b32e192f376a8cc325e7) feat(*): add data containers (@bacongobbler)
- [ef67e1b](https://github.com/deis/deis/commit/ef67e1be3eba4a0300308bf65d4a3362673a7b7a) docs(scheduler): removed most references to old concepts (@mboersma)
- [eaa2688](https://github.com/deis/deis/commit/eaa26886e84ca2bba1748c21bff00a7055a9a8d2) docs(scheduler): cleaned up index and concepts (@mboersma)
- [7be28e0](https://github.com/deis/deis/commit/7be28e0420c44ca01a3fc4aadbd14990b2f19072) docs(scheduler): added copy of draft architectural diagram (@mboersma)
- [36e578d](https://github.com/deis/deis/commit/36e578de4d8aebbf35f11cd65c228cf4b6c222e0) fix(*): use host's IP address (@bacongobbler)
- [d5aa3dd](https://github.com/deis/deis/commit/d5aa3ddc3a32902f62a66be6745ebaf8a657756b) docs(scheduler): add and update README.md files (@mboersma)
- [16f7621](https://github.com/deis/deis/commit/16f762101730c5ecbfd3baf3ffe807f63129a87b) docs(contrib/digitalocean): remove support for Digital Ocean (@carmstrong)
- [859940d](https://github.com/deis/deis/commit/859940d271e702f42ef28dde98d044a5064c92f5) fix(systemd): lookup default gateway interface for HOST_IP (@gabrtv)
- [4a1beb1](https://github.com/deis/deis/commit/4a1beb1d67d5aab87d293909634ae5aaee772cb9) fix(docker): add docker-patch.service for symlink workaround (@gabrtv)
- [41489c1](https://github.com/deis/deis/commit/41489c1d5a1110c66352be64102d715cf17f6bbe) fix(systemd): remove path to docker dev binary (@gabrtv)
- [df11de0](https://github.com/deis/deis/commit/df11de09ee4018ae592d66397e8a0bd6adf6ba1a) fix(builder): fix data container volume mounts (@gabrtv)
- [71fb3e6](https://github.com/deis/deis/commit/71fb3e688a86dc50429bb949f24a84a9a73a20f5) chore(builder): remove docker push debug prints (@gabrtv)
- [bb87d3c](https://github.com/deis/deis/commit/bb87d3c197e38215a29e2acf1c8b00486ad958ea) docs(contrib/ec2): add ssh-add command (@bacongobbler)
- [b8df76a](https://github.com/deis/deis/commit/b8df76a845236c3960bf86581b38f43f51505659) fix(registry): wait for port to become available in ExecPostStart (@gabrtv)
- [7dc4a3a](https://github.com/deis/deis/commit/7dc4a3a4f5ec369801e8f1696b7c4cf7aa0e731d) fix(user-data): support dynamic interface for seed-deis-registry (@gabrtv)
- [101c88b](https://github.com/deis/deis/commit/101c88b72bb0ebedbdd10fc1f033850f4fd801f0) fix(*): be more permissive with etcd (@bacongobbler)
- [1dcf378](https://github.com/deis/deis/commit/1dcf3786f110942e1562dcae722013793ff29bd7) chore(client): remove apps:calculate (@bacongobbler)
- [ecf167d](https://github.com/deis/deis/commit/ecf167d762d3153b9385a7c2faecdb3d727c2bad) fix(builder): reside beside the controller (@bacongobbler)
- [e3e8c18](https://github.com/deis/deis/commit/e3e8c183d40a9f5958c9f4d4ce2abe0f90e867f2) fix(fsm): save state transitions as they happen (@gabrtv)
- [b531857](https://github.com/deis/deis/commit/b53185703deceb455500253f53778cbeced383c3) fix(seed-registry): fix typo in ExecStart (@gabrtv)
- [06f4c79](https://github.com/deis/deis/commit/06f4c794890447bcc66ed6bc5d91f46b6a13ad3b) docs(contrib/rackspace): update to reflect CoreOS (@carmstrong)
- [e3c06dc](https://github.com/deis/deis/commit/e3c06dca9770a8ea064b574e8179bf89f48fdba5) docs(contrib): fix ec2 and rackspace formatting (@carmstrong)
- [af7a177](https://github.com/deis/deis/commit/af7a1774c091597379d9e6b370ad7d37030368c7) refactor(contrib): remove check-deis-deps (@carmstrong)
- [78f05d4](https://github.com/deis/deis/commit/78f05d4fb533a363002115cf8c8fe168b88d05ae) fix(app): restored apps:rollback functionality (@mboersma)
- [08d6d24](https://github.com/deis/deis/commit/08d6d240811434f7d9a5fc558650403c71fa5152) fix(contrib/ec2): add custom template (@bacongobbler)
- [6575e4f](https://github.com/deis/deis/commit/6575e4f27c9a1a9ca06f3842b7fd0e8a46e5a6cc) fix(contrib/ec2): be more restrictive with ports (@bacongobbler)
- [8965c0f](https://github.com/deis/deis/commit/8965c0f6c68268bc2c2976029fcbd14c4b699036) fix(contrib/ec2): add UDP perms to the logger (@bacongobbler)
- [b8ced1b](https://github.com/deis/deis/commit/b8ced1b5f1db32bf7c00d404ca0dba0977a95cd4) docs(contrib/ec2): use public DNS name (@bacongobbler)
- [7701372](https://github.com/deis/deis/commit/7701372aa44792ee53311a6ba5dc1d9f4a54d900) feat(logger): add new syslog daemon based on https://github.com/ziutek/syslog (@gabrtv)
- [7216f73](https://github.com/deis/deis/commit/7216f73ea4842a6387558b496347159f54de9e50) feat(logger): use new logger component from controller (@gabrtv)
- [4c99c2b](https://github.com/deis/deis/commit/4c99c2bc3bcd322a6092bdc04066f79fa3cb0037) fix(coreos): add wait_for_announcer to block until containers are ready (@gabrtv)
- [66f4b3d](https://github.com/deis/deis/commit/66f4b3d46e6a5161d19d72f67a9a188219770f9a) docs(scheduler): updated CLI instructions and developer docs (@mboersma)
- [4757378](https://github.com/deis/deis/commit/475737883cd3f1852e3606af65a24742bf8aa52d) fix(controller): use Deis' fork of django-fsm (@bacongobbler)
- [385dd6a](https://github.com/deis/deis/commit/385dd6ae9ff5ee05dc96e8507dec4da8fb46507b) docs(scheduler): add migration doc, update faq and localdev (@mboersma)
- [c35a25f](https://github.com/deis/deis/commit/c35a25f297f96e9a6b9708ebd58fde5274a45f19) fix(registry): Use cache host and port from confd (@johanneswuerbach)
- [1cd651e](https://github.com/deis/deis/commit/1cd651e1a2d88ee29a95b8f43dbc02bb8a926408) docs(scheduler): update "provision a controller" doc (@mboersma)
- [c264205](https://github.com/deis/deis/commit/c264205242541ae8c964e700669a39ec8321904c) docs(scheduler): mention coreos-vagrant reboot known issue (@mboersma)
- [ed6a3b8](https://github.com/deis/deis/commit/ed6a3b88c0474289553d0818e17459212d4f4560) fix(controller): use django.conf (@bacongobbler)
- [f45919c](https://github.com/deis/deis/commit/f45919c0c68fa09fd22492008013f64e137162cb) docs(scheduler): add reminder about migration and v0.7.0 (@mboersma)
- [#713](https://github.com/deis/deis/pull/713) Merge pull request #713 from opdemand/scheduler (@opdemand)
- [d9b8578](https://github.com/deis/deis/commit/d9b8578b0fc4c155d0428818574cfa5eb0ada3bd) fix(router): docker build fails when apt archive is stale (@mboersma)
- [#733](https://github.com/deis/deis/pull/733) Merge pull request #733 from opdemand/fix-router-build (@opdemand)
- [5f8f775](https://github.com/deis/deis/commit/5f8f77556030f04fe92e7e1b8652f22eb7b43516) fix(controller/systemd): schedule on same machine as logger (@carmstrong)
- [#734](https://github.com/deis/deis/pull/734) Merge pull request #734 from opdemand/fix_controller_unit (@opdemand)
- [28458f2](https://github.com/deis/deis/commit/28458f2397f499e14e8e5645851c19fe53a32dc0) fix(docs): mention "easy-fix" tag in contributing docs (@mboersma)
- [dfd202b](https://github.com/deis/deis/commit/dfd202b6e43d4440ef767ed5298df0e561898e93) docs(migrating): remove "nothing changes" phrase (@mboersma)
- [8eb5a01](https://github.com/deis/deis/commit/8eb5a01e85dc9e52ac4210e34cd2fac9436fa59e) docs(contrib/vagrant): update location of fleetctl (@bacongobbler)
- [#735](https://github.com/deis/deis/pull/735) Merge pull request #735 from opdemand/easy-fix-contrib (@opdemand)
- [e08de54](https://github.com/deis/deis/commit/e08de5496e53831dd21c4361f9b70686960cb239) fix(travis-ci): remove defunct "scheduler" branch (@mboersma)
- [0695711](https://github.com/deis/deis/commit/06957116f216783d78f245e25b2807299b69040e) fix(controller/tests): run coverage with "--timid" flag (@mboersma)
- [aae8fcd](https://github.com/deis/deis/commit/aae8fcd13fe0d82cd4bc5e5499e4acb2d4bcaa8a) chore(contrib/vagrant): lower memory to 2GB (@carmstrong)
- [3dda9fb](https://github.com/deis/deis/commit/3dda9fbeab1c546b8f137ea22bb4492d7e6f9112) fix(controller): fix container admin view (@bacongobbler)
- [746e140](https://github.com/deis/deis/commit/746e140cbd89aa29d3a813c270611394851bfb80) feat(controller): make coreos the default cluster (@bacongobbler)
- [ac966be](https://github.com/deis/deis/commit/ac966be8261dd9c0b024a19a2a8d50df8fb02f27) feat(controller): add faulty cluster (@bacongobbler)
- [#742](https://github.com/deis/deis/pull/742) Merge pull request #742 from opdemand/postgresql-test-errors (@opdemand)
- [8517d7d](https://github.com/deis/deis/commit/8517d7dff3e57b839a3c1825843e24486a8c5a22) test(controller): add more container state tests (@bacongobbler)
- [1d05e0d](https://github.com/deis/deis/commit/1d05e0d18061ad2559ecf16dead104dffb2c632c) fix(controller): remove cluster from list_filter (@bacongobbler)
- [33eee91](https://github.com/deis/deis/commit/33eee91c916b00fa8e84bd977be016ca847f8f01) fix(controller): fix flake8 (@bacongobbler)
- [#745](https://github.com/deis/deis/pull/745) Merge pull request #745 from opdemand/better-scheduling-tests (@opdemand)
- [#744](https://github.com/deis/deis/pull/744) Merge pull request #744 from opdemand/vagrant_memory (@opdemand)
- [#728](https://github.com/deis/deis/pull/728) Merge pull request #728 from johanneswuerbach/fix-prod-registry (@johanneswuerbach)
- [084dc18](https://github.com/deis/deis/commit/084dc18c264285280e744b4669ee4806901343ec) fix(vagrant): enable patch-docker (@bacongobbler)
- [0f6761e](https://github.com/deis/deis/commit/0f6761e23debb211bcd5b81b623a0862672670cc) feat(client): add --no-remote option (@bacongobbler)
- [9614d3b](https://github.com/deis/deis/commit/9614d3b6c1560f2c2d15cd0dc8aadafa0890cfa8) fix(contrib/vagrant): prevent Vagrant reboots (@carmstrong)
- [c9d0feb](https://github.com/deis/deis/commit/c9d0febcc5ca1f4260e338358f3076f901c26159) chore(contrib/vagrant): update CoreOS box image (@carmstrong)
- [5dc2abe](https://github.com/deis/deis/commit/5dc2abe784f565fdb20bf7e1345d42d57e62be23) refactor(client): fix program flow of --no-remote (@bacongobbler)
- [#754](https://github.com/deis/deis/pull/754) Merge pull request #754 from opdemand/stop_vagrant_reboots (@opdemand)
- [#757](https://github.com/deis/deis/pull/757) Merge pull request #757 from opdemand/hotfix-no-remote (@opdemand)
- [#750](https://github.com/deis/deis/pull/750) Merge pull request #750 from opdemand/enable-patched-docker (@opdemand)
- [#756](https://github.com/deis/deis/pull/756) Merge pull request #756 from opdemand/update_virtualbox_image (@opdemand)
- [4db0fe4](https://github.com/deis/deis/commit/4db0fe40f8ba026e384e17ba1a7cc7ba69473dc9) docs(concepts): add back build/release/run docs (@bacongobbler)
- [a9685e3](https://github.com/deis/deis/commit/a9685e3dd4fbcf6a14767f2d60204e4e882dcdd7) revert aafe0dd813591254471db597d452d75791f9280f (@bacongobbler)
- [#758](https://github.com/deis/deis/pull/758) Merge pull request #758 from opdemand/build-release-run (@opdemand)
- [af3bf60](https://github.com/deis/deis/commit/af3bf602ca2c7e0665d21e045d66689b4c5e1ee9) Fix contrib/vagrant/Makefile to work on non bash-symlinked /bin/sh. (@jperville)
- [cff6b59](https://github.com/deis/deis/commit/cff6b59502c9d4cb95cb661204e3137fe5e17687) fix(controller): update Django admin fields (@mboersma)
- [#761](https://github.com/deis/deis/pull/761) Merge pull request #761 from jperville/makefile-bashism-fix (@jperville)
- [3d09675](https://github.com/deis/deis/commit/3d096758cd9919b33c24819e20115e79f90ae07e) fix(contrib/vagrant): fix vagrant n-node Makefile (@carmstrong)
- [#759](https://github.com/deis/deis/pull/759) Merge pull request #759 from opdemand/vagrant_3node_makefile (@opdemand)
- [2b23db5](https://github.com/deis/deis/commit/2b23db5583c8f08feb800d54d5d0e3ee063fb9f7) feat(registry): add s3 redirect support (@johanneswuerbach)
- [#764](https://github.com/deis/deis/pull/764) Merge pull request #764 from johanneswuerbach/registry-s3-redirect (@johanneswuerbach)
- [d857735](https://github.com/deis/deis/commit/d85773542a54a972f9eb9ffb8493d3e35a44bc03) fix(logger): fix publish loop (@bacongobbler)
- [29f503a](https://github.com/deis/deis/commit/29f503afba0d95a713782a37052f6b6d6299ce87) feat(contrib/ec2): add m3 instances to template (@bacongobbler)
- [67ebd8c](https://github.com/deis/deis/commit/67ebd8c529812b15d90a4afbbc265e5f1717c635) refactor(seed-deis-registry): move to .service file (@carmstrong)
- [#770](https://github.com/deis/deis/pull/770) Merge pull request #770 from opdemand/move_seed_registry (@opdemand)
- [8f48f64](https://github.com/deis/deis/commit/8f48f64ef727e528c01357c1dce3beb69a9f0a56) docs(contrib): maintainers merge their own PRs (@bacongobbler)
- [b657d18](https://github.com/deis/deis/commit/b657d18d43eeae124d430e7fc418da31d5ed5049) fix(vagrant): use start on make run (@bacongobbler)
- [#763](https://github.com/deis/deis/pull/763) Merge pull request #763 from opdemand/fix-django-admin (@opdemand)
- [6ce0371](https://github.com/deis/deis/commit/6ce03710ff921a4fef72c619a083f3bc8cd90ec1) fix(seed-deis-registry): move into deis-registry.service (@carmstrong)
- [#775](https://github.com/deis/deis/pull/775) Merge pull request #775 from opdemand/seed_registry_fix (@opdemand)
- [#771](https://github.com/deis/deis/pull/771) Merge pull request #771 from opdemand/more-contrib-standards (@opdemand)
- [#772](https://github.com/deis/deis/pull/772) Merge pull request #772 from opdemand/use-make-start (@opdemand)
- [#766](https://github.com/deis/deis/pull/766) Merge pull request #766 from opdemand/fix-logger-publish-loop (@opdemand)
- [90c2ead](https://github.com/deis/deis/commit/90c2ead381c57194fc91aefb79a839b53010b111) fix(contrib/ec2): bump size to m3.large (@bacongobbler)
- [faf2c92](https://github.com/deis/deis/commit/faf2c920a0422254acc74c22393fe3e7bbc3ac6e) fix(docker): remove patched docker, switch to deis/base data container (@gabrtv)
- [a43e0d6](https://github.com/deis/deis/commit/a43e0d6a238a679983de6b3306a27fce69f74602) feat(coreos): update vagrant box to coreos-291.0.0 (@mboersma)
- [#778](https://github.com/deis/deis/pull/778) Merge pull request #778 from opdemand/remove-patched-docker (@opdemand)
- [#777](https://github.com/deis/deis/pull/777) Merge pull request #777 from opdemand/coreos-289.0.0 (@opdemand)
- [90dab33](https://github.com/deis/deis/commit/90dab33613be96e618d1def677db852c21b47010) fix(contrib/vagrant): fix Makefile loop (@carmstrong)
- [aab675a](https://github.com/deis/deis/commit/aab675aeed62d5e42a76e085387e710b09c5f04f) fix(controller): fix registry API calls (@bacongobbler)
- [#782](https://github.com/deis/deis/pull/782) Merge pull request #782 from opdemand/fix-registry-errors (@opdemand)
- [#768](https://github.com/deis/deis/pull/768) Merge pull request #768 from opdemand/add-m3-instances (@opdemand)
- [#780](https://github.com/deis/deis/pull/780) Merge pull request #780 from opdemand/fix_makefile_loop (@opdemand)
- [0c80d53](https://github.com/deis/deis/commit/0c80d53444a5d58f1c783b600b8d992a116b6d0e) fix(logger): import syslog package from github (@bacongobbler)
- [825a0c6](https://github.com/deis/deis/commit/825a0c65c4e1a2265f43d036c34fa48728417c20) refactor(seed-deis-registry): remove seed-deis-registry dependency (@carmstrong)
- [#788](https://github.com/deis/deis/pull/788) Merge pull request #788 from opdemand/remove_registry_block (@opdemand)
- [a678b23](https://github.com/deis/deis/commit/a678b239df793680b073c20bef5338f26e16285a) fix(Makefile): alphabetize component names (@bacongobbler)
- [#789](https://github.com/deis/deis/pull/789) Merge pull request #789 from opdemand/import-syslog (@opdemand)
- [d4bb423](https://github.com/deis/deis/commit/d4bb423bf0febbf0b69a5c7b5ec24dc4ff810b75) fix(logger): install deps in proper location (@bacongobbler)
- [2731b3a](https://github.com/deis/deis/commit/2731b3a352201313d0f24fd087dc9e5f03375743) chore(controller): remove provider libraries from requirements.txt (@mboersma)
- [#793](https://github.com/deis/deis/pull/793) Merge pull request #793 from opdemand/fix-syslogd (@opdemand)
- [#790](https://github.com/deis/deis/pull/790) Merge pull request #790 from opdemand/alphabetized (@opdemand)
- [0e32e03](https://github.com/deis/deis/commit/0e32e039397730701980906bc5306f24e98e970c) chore(contrib/vagrant): update vagrant box to coreos-295.0.0 (@mboersma)
- [#795](https://github.com/deis/deis/pull/795) Merge pull request #795 from opdemand/requirements-diet (@opdemand)
- [#796](https://github.com/deis/deis/pull/796) Merge pull request #796 from deis/coreos-295.0.0 (@deis)
- [7ac2e85](https://github.com/deis/deis/commit/7ac2e858d3457aa544a9108743e1c2564998ecb2) fix(github): update URL references to opdemand/deis (@mboersma)
- [#801](https://github.com/deis/deis/pull/801) Merge pull request #801 from deis/deis-deis-fixup (@deis)
- [cd66fd2](https://github.com/deis/deis/commit/cd66fd281181384ec9c4f9e7c5cd18aa89da1923) refactor(Vagrantfile): merge Vagrantfiles and update to upstream (@carmstrong)
- [241c696](https://github.com/deis/deis/commit/241c6960c85ce329377ba5699c805b1ea19c6c80) fix(tests): skip problematic section of BuildTest (@mboersma)
- [3b2ee0d](https://github.com/deis/deis/commit/3b2ee0d7dd14178eb99bc609b396798aca527f6a) fix(coverage): revert addition of "--timid" flag to coverage (@mboersma)
- [#803](https://github.com/deis/deis/pull/803) Merge pull request #803 from deis/postgres-tests-workaround (@deis)
- [#799](https://github.com/deis/deis/pull/799) Merge pull request #799 from deis/merge_vagrantfiles (@deis)
- [828a1df](https://github.com/deis/deis/commit/828a1dfc22c07ac9d22419bacd4e33805b5f393f) fix(Makefile): upload fleetctl units all at once (@carmstrong)
- [71a86d9](https://github.com/deis/deis/commit/71a86d9d08ca654344dabf60bd61cfdfba29163c) fix(README): remove extraneous ` (@carmstrong)
- [#807](https://github.com/deis/deis/pull/807) Merge pull request #807 from deis/fix_fleetctl_scheduling (@deis)
- [2ce2c3e](https://github.com/deis/deis/commit/2ce2c3e10b851db9967ffdcb0b6563978eccf3d2) chore(contrib/vagrant): update vagrant box to coreos-296.0.0 (@mboersma)
- [#808](https://github.com/deis/deis/pull/808) Merge pull request #808 from deis/coreos-296.0.0 (@deis)
- [3ddbccf](https://github.com/deis/deis/commit/3ddbccfddf694c783274adf1d0a903403927a1b6) refactor(client): rename containers to ps (@bacongobbler)
- [0f5c460](https://github.com/deis/deis/commit/0f5c460c07dee0347b6c92b28d33cc37a8e59d7d) fix(scheduler): announce proper address (@bacongobbler)
- [#813](https://github.com/deis/deis/pull/813) Merge pull request #813 from deis/rename-containers (@deis)
- [#814](https://github.com/deis/deis/pull/814) Merge pull request #814 from deis/fix-announce (@deis)
- [123776e](https://github.com/deis/deis/commit/123776e9d1dd8224b6269f423ff40eed344d40b4) Removed duplicate etcdctl installation statements (@zyegfryed)
- [#822](https://github.com/deis/deis/pull/822) Merge pull request #822 from zyegfryed/patch-1 (@zyegfryed)
- [69f138a](https://github.com/deis/deis/commit/69f138a4cfc17c00222460454975ef799f031a7c) chore(vagrant): update box to coreos-298.0.0 (@mboersma)
- [#824](https://github.com/deis/deis/pull/824) Merge pull request #824 from deis/coreos-298.0.0 (@deis)
- [787d857](https://github.com/deis/deis/commit/787d857e97e46fed661bb0c51fc651cb8a3abc46) fix(services): manage dependencies manually (@carmstrong)
- [#821](https://github.com/deis/deis/pull/821) Merge pull request #821 from deis/fix_scheduling (@deis)
- [04b21e4](https://github.com/deis/deis/commit/04b21e4df88cfdbb5de5df5d7664751fa386591a) chore(controller): update to Django 1.6.3 security release (@mboersma)
- [6783d9a](https://github.com/deis/deis/commit/6783d9a787c719d98f032591d8a5851a20849dea) fix(contrib/ec2): fix ephemeral port range (@bacongobbler)
- [4a029f6](https://github.com/deis/deis/commit/4a029f63ed5accdeee934fad9509bd20d1cc5c59) fix(coverage): cd to controller dir before submitting to coveralls.io (@mboersma)
- [#827](https://github.com/deis/deis/pull/827) Merge pull request #827 from deis/fix-ephemeral-port-range (@deis)
- [#828](https://github.com/deis/deis/pull/828) Merge pull request #828 from deis/fix-coveralls-io (@deis)
- [#794](https://github.com/deis/deis/pull/794) Merge pull request #794 from deis/django-1.6.3 (@deis)
- [1fb49d8](https://github.com/deis/deis/commit/1fb49d8493839fb69eee5247d3b6b65ea1f9c3f4) fix(Makefile): use /bin/sh if syntax (@bacongobbler)
- [#829](https://github.com/deis/deis/pull/829) Merge pull request #829 from deis/make-controller (@deis)
- [a2d006c](https://github.com/deis/deis/commit/a2d006cd8b8a5c9336319bee98835d7b54610145) fix(builder): update Docker to v0.10.0 (@mboersma)
- [f112a5f](https://github.com/deis/deis/commit/f112a5f0e0d2eaa2bdbdecf5fe2b7f084f2d4179) feat(Makefile): pull deis/slugrunner (@bacongobbler)
- [#832](https://github.com/deis/deis/pull/832) Merge pull request #832 from deis/slugrunner-pull (@deis)
- [b81f4a8](https://github.com/deis/deis/commit/b81f4a85abc2ea3980d998a0bebde812e4c53ae4) feat(controller): make Procfile mandatory (@bacongobbler)
- [a4da57d](https://github.com/deis/deis/commit/a4da57daac768e0f5e49f32d6c714ffcdcd6d915) fix(Makefile): 'make run' starts Deis in one shot (@carmstrong)
- [#833](https://github.com/deis/deis/pull/833) Merge pull request #833 from deis/builder-docker-0.10.0 (@deis)
- [#831](https://github.com/deis/deis/pull/831) Merge pull request #831 from deis/makefile_friendly_loop (@deis)
- [c7a3247](https://github.com/deis/deis/commit/c7a3247fafd086554516991d1131da0c4ed9b05b) chore(contrib/ec2): update AMIs for build 298.0.0 (@bacongobbler)
- [326d383](https://github.com/deis/deis/commit/326d3838d72f3752609d6a50c334e11a7cf377b5) chore(contrib): remove docker project (@bacongobbler)
- [#835](https://github.com/deis/deis/pull/835) Merge pull request #835 from deis/update-amis (@deis)
- [#836](https://github.com/deis/deis/pull/836) Merge pull request #836 from deis/remove-docker-contrib (@deis)
- [bc5f72b](https://github.com/deis/deis/commit/bc5f72bf2cf09e3a31622aed45fcbb286dee5a5a) fix(controller): rename destroy() to delete() and call its super() (@mboersma)
- [027a001](https://github.com/deis/deis/commit/027a0019c410cdb88d6cbb6a85a75344be6c7d02) chore(controller): update to Django 1.6.4 bugfix release (@mboersma)
- [c90614a](https://github.com/deis/deis/commit/c90614adecaa31fe7a865fa9dc42423a3fece1c1) fix(docs): rename client/containers.rst to ps.rst (@mboersma)
- [42d2122](https://github.com/deis/deis/commit/42d2122b978fece5e87894fae8b68b42a26139f6) docs(localdev): add "Test Your Changes" section to local dev docs (@mboersma)
- [2a0bd99](https://github.com/deis/deis/commit/2a0bd99c37572b3f8c34bbf216c7730a538c8c90) docs(ec2): remove reference to obsolete script (@mboersma)
- [#842](https://github.com/deis/deis/pull/842) Merge pull request #842 from deis/add-interactive-dev-docs (@deis)
- [#841](https://github.com/deis/deis/pull/841) Merge pull request #841 from deis/fix-client-ps-docs (@deis)
- [#840](https://github.com/deis/deis/pull/840) Merge pull request #840 from deis/django-1.6.4 (@deis)
- [#839](https://github.com/deis/deis/pull/839) Merge pull request #839 from deis/app-fleet-delete (@deis)
- [#844](https://github.com/deis/deis/pull/844) Merge pull request #844 from deis/ec2-readme-cleanup (@deis)
- [#834](https://github.com/deis/deis/pull/834) Merge pull request #834 from deis/procfile-mandatory (@deis)
- [c2f4085](https://github.com/deis/deis/commit/c2f4085ed5f26c46ccc1c07cf8c289d2aea78837) feat(test): add integration tests (@carmstrong)
- [#837](https://github.com/deis/deis/pull/837) Merge pull request #837 from deis/integration_tests (@deis)
- [ebf51c2](https://github.com/deis/deis/commit/ebf51c23b8d903385714355571e6e05e8ec256da) docs(contrib/ec2): generate the keypair (@bacongobbler)
- [ac0ae39](https://github.com/deis/deis/commit/ac0ae39fdae724e64cca7c7bc8181b1382286d6b) docs(contrib/ec2): use ~/.ssh/deis instead (@bacongobbler)
- [#853](https://github.com/deis/deis/pull/853) Merge pull request #853 from deis/ec2-keypair (@deis)
- [df5f1d1](https://github.com/deis/deis/commit/df5f1d13af6beb0fc888ae9a474de9d2ef2c460a) feat(Makefile): status reports and better error handling (@carmstrong)
- [50b0cda](https://github.com/deis/deis/commit/50b0cdab0c6a6fd7741bd0a12e8dc784694c39bd) docs(release): update procedure for v0.8.0 release (@mboersma)
- [#854](https://github.com/deis/deis/pull/854) Merge pull request #854 from deis/more_makefile_error_checking (@deis)
- [#855](https://github.com/deis/deis/pull/855) Merge pull request #855 from deis/update-release-doc (@deis)
- [2c0e4c6](https://github.com/deis/deis/commit/2c0e4c60194a9f736c0a486c2285d1c91b7fef8a) fix(services): use COREOS_PUBLIC_IPV4 for advertised IP (@carmstrong)
- [#851](https://github.com/deis/deis/pull/851) Merge pull request #851 from deis/dont_guess_advertised_ip (@deis)
- [8d2f30a](https://github.com/deis/deis/commit/8d2f30a715f21c40f51230266646a00427a83ec4) docs(client): update links to binary versions for v0.8.0 (@mboersma)
- [b493ad9](https://github.com/deis/deis/commit/b493ad9961480b2635491c2e1a7d2b845715549e) fix(builder): start after controller, block until running (@carmstrong)
- [4f92caf](https://github.com/deis/deis/commit/4f92caf601e161c72ef687c5f3ca5eb68d6226a5) Revert "feat(controller): make Procfile mandatory" (@bacongobbler)
- [d1e14a6](https://github.com/deis/deis/commit/d1e14a6f76789a86f1959ba7bf6057b3abe3b960) feat(controller): handle special cmd proctype (@bacongobbler)
- [#861](https://github.com/deis/deis/pull/861) Merge pull request #861 from deis/binary-clients-v0.8.0 (@deis)
- [#862](https://github.com/deis/deis/pull/862) Merge pull request #862 from deis/cmd-proctype (@deis)
- [#863](https://github.com/deis/deis/pull/863) Merge pull request #863 from deis/block_builder (@deis)
- [198e550](https://github.com/deis/deis/commit/198e55005faf45ab0491bfa07d069a87a07a030c) docs(readme): move cluster testing section out of developer workflow (@gabrtv)
- [8cf66bd](https://github.com/deis/deis/commit/8cf66bd3328a8185096dff59be5fbc6470fde92d) docs(readme): add ref to integration tests, reorder language (@gabrtv)
- [#864](https://github.com/deis/deis/pull/864) Merge pull request #864 from deis/reorder-readme (@deis)
- [1286afc](https://github.com/deis/deis/commit/1286afc86f5d45abafc1e896e5784a39030595d4) docs(ec2): add CLI key import, explicit chdir (@gabrtv)
- [a8869a9](https://github.com/deis/deis/commit/a8869a98d11ca820ccd822f5dbc55810356d58ed) docs(rackspace): add explicit chdir (@gabrtv)
- [cbeca5a](https://github.com/deis/deis/commit/cbeca5ae4fd3726619ec0c5ca6bf24fec53e8c91) docs(ec2): add CLI key import, explicit chdir (@gabrtv)
- [#865](https://github.com/deis/deis/pull/865) Merge pull request #865 from deis/fix-ec2-readme (@deis)
- [57ceade](https://github.com/deis/deis/commit/57ceade32933c583b80ad0888c2df96e81109676) fix(contrib/ec2): fix README from #865 (@carmstrong)
- [#868](https://github.com/deis/deis/pull/868) Merge pull request #868 from deis/fix_865 (@deis)
- [5d3cebf](https://github.com/deis/deis/commit/5d3cebf9e29654b70790585d404aa1189cb05b56) refactor(Makefile): better multi-node support (@carmstrong)
- [66660a4](https://github.com/deis/deis/commit/66660a41b8be0a1e795a77973394737fed0f1044) fix(systemd): handle starting containers that already exist (@gabrtv)
- [#871](https://github.com/deis/deis/pull/871) Merge pull request #871 from deis/remove-on-start (@deis)
- [0426805](https://github.com/deis/deis/commit/04268059cc2f9e76814c742cda417a8d1bd7bb7e) docs(dns): move content from README.md into "Configure DNS" doc (@mboersma)
- [#870](https://github.com/deis/deis/pull/870) Merge pull request #870 from deis/dns_readme (@deis)
- [09dc820](https://github.com/deis/deis/commit/09dc8207a0d34f2b6d124e4bba2798fab50eb2bc) fix(docker): add an ARP entry each time a container comes up (@mboersma)
- [9502968](https://github.com/deis/deis/commit/9502968080cc07e43e37fc71fbe13ebc192772ae) fix(contrib): #859 remove update manager stop commands from Vagrantfile (@dginther)
- [8dfbe8c](https://github.com/deis/deis/commit/8dfbe8c5ce6aac0f62c695ec722fe38c841b1816) fix(user-data): use $private_ipv4 for etcd/fleet (@carmstrong)
- [9e3b471](https://github.com/deis/deis/commit/9e3b4712a16d08828aa2c8ec3d4aad100cb9895c) fix(contrib/ec2): use shared user-data (@carmstrong)
- [a99c75a](https://github.com/deis/deis/commit/a99c75ae8b222eda0843e036f9d06fef5ff0127f) docs(typo): change "EC2" to "Rackspace" (@mboersma)
- [#869](https://github.com/deis/deis/pull/869) Merge pull request #869 from deis/contrib_use_user_data (@deis)
- [0818b79](https://github.com/deis/deis/commit/0818b79484c78bc321c9e8381ed0283358a4757c) style(coreos): remove unused systemd template strings (@mboersma)
- [22681a3](https://github.com/deis/deis/commit/22681a3aa6b2f27ac65a1033040104754d54da23) fix(controller): start logger service after container (@mboersma)
- [c383ba1](https://github.com/deis/deis/commit/c383ba1ef8bb376edcdc965bc22a97127422c773) refactor(controller): remove duplicate run() method (@mboersma)
- [#877](https://github.com/deis/deis/pull/877) Merge pull request #877 from deis/logger-after-container (@deis)
- [#873](https://github.com/deis/deis/pull/873) Merge pull request #873 from deis/prime-arp-workaround (@deis)
- [aac4f6a](https://github.com/deis/deis/commit/aac4f6a5198e8060032d42a9f411167f8c0a51f7) fix(Makefile): check for errors after controller (@carmstrong)
- [#881](https://github.com/deis/deis/pull/881) Merge pull request #881 from deis/add_check_errors (@deis)
- [#879](https://github.com/deis/deis/pull/879) Merge pull request #879 from deis/vestigial-templates (@deis)
- [5521435](https://github.com/deis/deis/commit/5521435c8a82dd555fea0b7330ca7b17af092a37) fix(systemd): handle app containers on restarts (@bacongobbler)
- [#872](https://github.com/deis/deis/pull/872) Merge pull request #872 from deis/remove-containers-on-start (@deis)
- [dd6e16d](https://github.com/deis/deis/commit/dd6e16deefd5f435aebcc2ce9f556b0716a876e5) docs(contributing): add skip ci tag (@bacongobbler)
- [fb38eb6](https://github.com/deis/deis/commit/fb38eb65ab7ad50386216f470e1656dbee3b0cd2) fix(test): specify machine IPs (@carmstrong)
- [#886](https://github.com/deis/deis/pull/886) Merge pull request #886 from deis/docs-no-ci (@deis)
- [ff24d62](https://github.com/deis/deis/commit/ff24d62c43f0724c7ba47e13bc3cbb48809bec04) fix(ec2): add explicit 100GB root volume (@gabrtv)
- [#889](https://github.com/deis/deis/pull/889) Merge pull request #889 from deis/test_hosts (@deis)
- [#895](https://github.com/deis/deis/pull/895) Merge pull request #895 from deis/fix-ec2-space (@deis)

### v0.7.0 (2014/04/10 22:42 +00:00)
- [273551c](https://github.com/deis/deis/commit/273551cd7cb65a329156385966cccc127a213a27) Switch master to v0.7.0. (@mboersma)
- [a6a3410](https://github.com/deis/deis/commit/a6a3410625da298fd0594fa1797c1268183844f8) Standardize containers on WORKDIR /app (@mboersma)
- [5e672cd](https://github.com/deis/deis/commit/5e672cda1e03aecd75e3e9d729d18302bcbb9067) Standardize install of pip 1.5.4. (@mboersma)
- [ac12f9a](https://github.com/deis/deis/commit/ac12f9abdc6ff2de3b9e40feaa7ce7faadf827d2) Clarify binding to $PORT with Dockerfiles (@gabrtv)
- [c20082c](https://github.com/deis/deis/commit/c20082c1ef5b074cf5d5b869675799b962301e81) Restored deis data bags creation to provision script. (@mboersma)
- [#665](https://github.com/deis/deis/pull/665) Merge pull request #665 from opdemand/664-vagrant-data-bags (@opdemand)
- [31af241](https://github.com/deis/deis/commit/31af241b2bf767311f3afff938aaf57ad52422b4) Clarify 2-LGTM code approval policy. (@mboersma)
- [#666](https://github.com/deis/deis/pull/666) Merge pull request #666 from opdemand/merge-approval (@opdemand)
- [a1e2ee0](https://github.com/deis/deis/commit/a1e2ee07e25ddde294d41c8c008b37521074f9c7) Adds friendly name in VirtualBox and Vagrant (@carmstrong)
- [#663](https://github.com/deis/deis/pull/663) Merge pull request #663 from opdemand/dockerfile-makefile-cleanup (@opdemand)
- [060cfc1](https://github.com/deis/deis/commit/060cfc1f1d55e45d681f8eba729261ca17f2e44a) Let berkshelf resolve dependencies itself (@carmstrong)
- [38647cb](https://github.com/deis/deis/commit/38647cbfd5dc79b38b74aa3cb3ffa611a8122aa7) pin pep8, pyflakes and flake8 to avoid travis errors (@gabrtv)
- [#670](https://github.com/deis/deis/pull/670) Merge pull request #670 from opdemand/pin-pep-pyflakes (@opdemand)
- [#669](https://github.com/deis/deis/pull/669) Merge pull request #669 from opdemand/carmstrong/fix_berksfile (@opdemand)
- [#667](https://github.com/deis/deis/pull/667) Merge pull request #667 from opdemand/carmstrong/vagrant_name (@opdemand)
- [#651](https://github.com/deis/deis/pull/651) Merge pull request #651 from opdemand/548-update-docker (@opdemand)
- [6a2ca5e](https://github.com/deis/deis/commit/6a2ca5e41adaaa2f66e05087bd039c0b04011031) Install lxc explicitly (@carmstrong)
- [d749d14](https://github.com/deis/deis/commit/d749d14190dcde1b18bb3d45db1b517e3af80194) Added major *scheduler* branch to Travis CI. (@mboersma)
- [5a5392f](https://github.com/deis/deis/commit/5a5392f4ee5cacfa2715f2a864adfae0ece62dd0) Reimplemented `deis shortcuts` command. (@mboersma)
- [#676](https://github.com/deis/deis/pull/676) Merge pull request #676 from opdemand/travis-scheduler-branch (@opdemand)
- [#677](https://github.com/deis/deis/pull/677) Merge pull request #677 from opdemand/deis-shortcuts (@opdemand)
- [deb6a14](https://github.com/deis/deis/commit/deb6a14c02a0fae2990b95f7f0d3c84b83984c29) dynamically disable registration (@bacongobbler)
- [3c6d792](https://github.com/deis/deis/commit/3c6d7925aebea9bdf3e66f7b586126cb5f1d90d7) import from django.conf instead (@bacongobbler)
- [d66bc17](https://github.com/deis/deis/commit/d66bc178c09fef1274649d7bdc59ad6fd00e8e2e) docs(standards): add commit style guide (@bacongobbler)
- [88a32c5](https://github.com/deis/deis/commit/88a32c5b425cdc71963adc83a33b1b533376fc60) docs(standards): update old commands (@bacongobbler)
- [6833b28](https://github.com/deis/deis/commit/6833b28603f2846ec52498a20c50905419ba146a) docs(standards): use proper english (@bacongobbler)
- [#672](https://github.com/deis/deis/pull/672) Merge pull request #672 from opdemand/carmstrong/install_lxc (@opdemand)
- [e99cc41](https://github.com/deis/deis/commit/e99cc41a3fd07f46d60f53ed05e26f7b39ae4a88) fix(logger): update rsyslog repo endpoint (@bacongobbler)
- [#681](https://github.com/deis/deis/pull/681) Merge pull request #681 from opdemand/rsyslog-repo-update (@opdemand)
- [207d7a7](https://github.com/deis/deis/commit/207d7a71970ad8669a970a7f3f7dc240c0ab55b9) docs(chef): update for the new manage.opscode.com (@mboersma)
- [#682](https://github.com/deis/deis/pull/682) Merge pull request #682 from opdemand/update-chef-admins (@opdemand)
- [#679](https://github.com/deis/deis/pull/679) Merge pull request #679 from opdemand/disable-registration (@opdemand)
- [#653](https://github.com/deis/deis/pull/653) Merge pull request #653 from opdemand/makefile (@opdemand)
- [#680](https://github.com/deis/deis/pull/680) Merge pull request #680 from opdemand/add-contributing (@opdemand)
- [bdb9f7c](https://github.com/deis/deis/commit/bdb9f7c3bc0d922e751bec372bcd40d926500a52) fix(docs): add codeblock spacing (@bacongobbler)
- [#683](https://github.com/deis/deis/pull/683) Merge pull request #683 from opdemand/fix-docs-formatting (@opdemand)
- [2b19dfd](https://github.com/deis/deis/commit/2b19dfd6911085b2b8db6d50d7a6f873196f8461) use force yes to install required packages (@jstop)
- [#688](https://github.com/deis/deis/pull/688) Merge pull request #688 from clyphub/master (@clyphub)
- [bd5f76f](https://github.com/deis/deis/commit/bd5f76f364c451f0f27fdd1a1821b410bf625226) fix(providers): update dop calls (@bacongobbler)
- [#691](https://github.com/deis/deis/pull/691) Merge pull request #691 from opdemand/fix-digitalocean (@opdemand)
- [#690](https://github.com/deis/deis/pull/690) Merge pull request #690 from opdemand/carmstrong/contributing_docs_stove (@opdemand)
- [#693](https://github.com/deis/deis/pull/693) Merge pull request #693 from opdemand/fix-logger (@opdemand)
- [d2dacb2](https://github.com/deis/deis/commit/d2dacb282c1fca7b104201f0b963ebe032f1e8e4) docs(components): explain a container's filesystem (@bacongobbler)
- [49b2907](https://github.com/deis/deis/commit/49b2907fd02687ff455028a261b7056220c77e6f) Updated CHANGELOG.md. (@mboersma)

### v0.6.0 (2014/03/24 02:55 +00:00)
- [d741c5f](https://github.com/deis/deis/commit/d741c5f52325ad203ca601ac3208be357f96c0ff) Switch master to v0.6.0. (@mboersma)
- [a51d0d2](https://github.com/deis/deis/commit/a51d0d2fce3d429b93b998d34840a77d4f07fc72) import submodules into project (@bacongobbler)
- [60da28e](https://github.com/deis/deis/commit/60da28ea49cd22a1160f172f5e752dff03f832d1) remove unused scripts (@bacongobbler)
- [822762f](https://github.com/deis/deis/commit/822762f23faf545362621e1625032324593766bf) move deis controller to separate project (@bacongobbler)
- [56603dc](https://github.com/deis/deis/commit/56603dcae81c70b42012e12761f8f2fc677302be) updated test suite to point to controller project (@bacongobbler)
- [887732e](https://github.com/deis/deis/commit/887732e92e2641f6b6794f3a0b29fbd221e5efb4) remove client module from controller settings (@bacongobbler)
- [f9f60b3](https://github.com/deis/deis/commit/f9f60b3df7b078d2b23c97551853141200f7991d) fix docs path issue (@bacongobbler)
- [3d079b9](https://github.com/deis/deis/commit/3d079b9e7f62a26309f7e7ee2ab91958619d0fa6) import deis/server (@bacongobbler)
- [7595ecf](https://github.com/deis/deis/commit/7595ecffd3e15e49eb188f9f091acd33bb1399b6) fix test locations (@bacongobbler)
- [bab9616](https://github.com/deis/deis/commit/bab96161998cba2d20bf066b2c9014046886b9a5) set exec bit on controller/bin/boot (@bacongobbler)
- [45c78fb](https://github.com/deis/deis/commit/45c78fb5fc4985cf9fa522cdf257dccfc5c4090a) use the -C option instead (@bacongobbler)
- [3b9518e](https://github.com/deis/deis/commit/3b9518ecfd82dffa4e2c0ceeb4e0c157aa7a7630) move `make client_binary` to client project (@bacongobbler)
- [3310a23](https://github.com/deis/deis/commit/3310a23d396af35218be68a8967c9903c7ba88c2) add back init (@bacongobbler)
- [44eab43](https://github.com/deis/deis/commit/44eab43cef9eadf3226673ec7f4dc22aaa2250c6) Fixed PYTHONPATH for Sphinx docs generation. (@mboersma)
- [b5440da](https://github.com/deis/deis/commit/b5440daa79b132adf7d011d0d33318de80698f92) Install latest pip. (@mboersma)
- [#641](https://github.com/deis/deis/pull/641) Merge pull request #641 from opdemand/godmode (@opdemand)
- [2a63d16](https://github.com/deis/deis/commit/2a63d16c1857eb5727e79f891ec75c45c538bfc4) add build command (@bacongobbler)
- [ebc5acb](https://github.com/deis/deis/commit/ebc5acbc89416ce93f0473f8500658770a47b1b6) Pin pip at v1.5.4. (@mboersma)
- [#644](https://github.com/deis/deis/pull/644) Merge pull request #644 from opdemand/642-pin-pip (@opdemand)
- [884fed5](https://github.com/deis/deis/commit/884fed539dd28d4174d7ee640ac391565bcecbdc) Replaced refs to server and worker with controller. (@mboersma)
- [2931ac6](https://github.com/deis/deis/commit/2931ac6cb3d60aeb0d1b802c7d545e371dc5b56c) fix confd wait issue (@bacongobbler)
- [6fe22d9](https://github.com/deis/deis/commit/6fe22d9429146729f2f3e0e441b45f9f1e5f03e4) only write out chef config if necessary, restart based on pid files (@gabrtv)
- [d179eed](https://github.com/deis/deis/commit/d179eedbd775e8947da455e09fcb2fb6fcc7feac) fix overlapping bind-mount issues (@bacongobbler)
- [#646](https://github.com/deis/deis/pull/646) Merge pull request #646 from opdemand/fix-registry-seed (@opdemand)
- [#635](https://github.com/deis/deis/pull/635) Merge pull request #635 from opdemand/pin-ruby-version (@opdemand)
- [#643](https://github.com/deis/deis/pull/643) Merge pull request #643 from opdemand/add-build-cmd (@opdemand)
- [#645](https://github.com/deis/deis/pull/645) Merge pull request #645 from opdemand/deis-controller-refs (@opdemand)
- [ce5b1bf](https://github.com/deis/deis/commit/ce5b1bfd8dd34bfc014806181befc5be4e1fb107) Run django unit tests via 'make test'. (@mboersma)
- [#650](https://github.com/deis/deis/pull/650) Merge pull request #650 from opdemand/fix-controller-make-test (@opdemand)
- [451711b](https://github.com/deis/deis/commit/451711bd41cc12eed36ecdb9aee252bfbb4c5e25) bump docker to v0.9.0 (@bacongobbler)
- [ba235f4](https://github.com/deis/deis/commit/ba235f4def550f2e7b760df5b4823d44ebbf4cb5) Added release instructions to create CHANGELOG.md files. (@mboersma)
- [063a328](https://github.com/deis/deis/commit/063a328f2b6484910709e16ed6b6a5425b84cd88) let each project manage how they build (@bacongobbler)
- [18f8d38](https://github.com/deis/deis/commit/18f8d38e1eb978fe16c0301ea6906f44a225f36e) Hide 'Versions' links unless at readthedocs.org, fixes #602. (@mboersma)
- [#654](https://github.com/deis/deis/pull/654) Merge pull request #654 from opdemand/fix-docs-versions-link (@opdemand)
- [#652](https://github.com/deis/deis/pull/652) Merge pull request #652 from opdemand/changelog-md (@opdemand)
- [7aa925c](https://github.com/deis/deis/commit/7aa925caf287280440088ad6f19826063453643c) Updated CLI binaries and links in client README. (@mboersma)
- [#658](https://github.com/deis/deis/pull/658) Merge pull request #658 from opdemand/update-cli-binaries (@opdemand)
- [bac1bcd](https://github.com/deis/deis/commit/bac1bcd4454178dcd1d88e76757da292dde5a387) Updated EC2 AMIs for v0.6.0 release. (@mboersma)
- [#659](https://github.com/deis/deis/pull/659) Merge pull request #659 from opdemand/time-to-make-the-amis (@opdemand)
- [7be8b3c](https://github.com/deis/deis/commit/7be8b3c5fbe45589cb1030c677f775f788aabc95) make SSH key and host nodes dir available (@bacongobbler)
- [99a43de](https://github.com/deis/deis/commit/99a43defd92f3f4cbab30a3a0a08c656622a11cc) add empty host_nodes_dir (@bacongobbler)
- [b55f6ff](https://github.com/deis/deis/commit/b55f6ff648da309d414b3d45d35b62ce6dbcde78) fix typo (@bacongobbler)
- [6550851](https://github.com/deis/deis/commit/65508513d4d78de7d94edcc9d5ad279100434b67) Revert "add empty host_nodes_dir" (@bacongobbler)
- [ffafee4](https://github.com/deis/deis/commit/ffafee4b9d9a25c335b01c31385bf6139bdb0dbc) fix path issue (@bacongobbler)
- [3f58af1](https://github.com/deis/deis/commit/3f58af1e94c3858f886cb4b719bd455933dbbf3a) fix dop images function call (@bacongobbler)
- [#662](https://github.com/deis/deis/pull/662) Merge pull request #662 from opdemand/fix-dop-images (@opdemand)
- [#661](https://github.com/deis/deis/pull/661) Merge pull request #661 from opdemand/655-vagrant-nodes (@opdemand)
- [f05ea96](https://github.com/deis/deis/commit/f05ea963779c3b5eb230e537b6b366558069c898) Updated CHANGELOG.md. (@mboersma)

### v0.5.2 (2014/03/18 15:42 +00:00)
- [1f8e06f](https://github.com/deis/deis/commit/1f8e06f963d20f4b51b91824e976da4433d930f9) Switch master to v0.5.2. (@mboersma)
- [f039d79](https://github.com/deis/deis/commit/f039d79703b5f4fa3565dce31666db40bd5c3f50) update script name (@bacongobbler)
- [#601](https://github.com/deis/deis/pull/601) Merge pull request #601 from opdemand/fix-rackspace-provider (@opdemand)
- [#609](https://github.com/deis/deis/pull/609) Merge pull request #609 from opdemand/vagrant-readme (@opdemand)
- [832b868](https://github.com/deis/deis/commit/832b86847316d93cdcd623a5818929ab612c9768) update developer docs to match containerize (@bacongobbler)
- [b0cc506](https://github.com/deis/deis/commit/b0cc506330e880e2e090074710d45f32956687e7) move etcd import to external module imports (@bacongobbler)
- [af4a9d5](https://github.com/deis/deis/commit/af4a9d5df89558e9000c2baf5b58fac20b8e8bf1) typo: use --formation option (@bacongobbler)
- [cf928d7](https://github.com/deis/deis/commit/cf928d7cd7441934ba3a953040c37cdead87169a) add Dockerfile docs and example apps (@bacongobbler)
- [79f1bfa](https://github.com/deis/deis/commit/79f1bfafd1fb6a27450e2a49e64d6554576558ab) elaborate on app deployment (@bacongobbler)
- [#611](https://github.com/deis/deis/pull/611) Merge pull request #611 from opdemand/534-dev-docs-update (@opdemand)
- [fef8241](https://github.com/deis/deis/commit/fef8241fc556ae5eb157c1f012b7574c0c43d255) add docs for buildpacks and dockerfiles (@bacongobbler)
- [44a040f](https://github.com/deis/deis/commit/44a040f5598b1c4aba5c83ec19efad081563dfff) just use links instead of the URL (@bacongobbler)
- [0745b27](https://github.com/deis/deis/commit/0745b275038214c0c1c4784716bf5070c1136dff) update local development docs (@bacongobbler)
- [3b91ff4](https://github.com/deis/deis/commit/3b91ff4d87f95be9e3986533303f692256ed0b7b) update buildstep to slugbuilder (@bacongobbler)
- [0e71835](https://github.com/deis/deis/commit/0e718357c38cc9341df3fa15a66908d8ab75ec21) remove double spaces (@bacongobbler)
- [72363e7](https://github.com/deis/deis/commit/72363e794029ffec562ba51ef7f0cad4a00c0cd5) update Deis' architecture docs (@bacongobbler)
- [3d12f6e](https://github.com/deis/deis/commit/3d12f6ef375c8df343716dafe856b2fcff877ff1) update operator's documentation (@bacongobbler)
- [#615](https://github.com/deis/deis/pull/615) Merge pull request #615 from opdemand/operations-update-docs (@opdemand)
- [#612](https://github.com/deis/deis/pull/612) Merge pull request #612 from opdemand/604-dockerfile-docs (@opdemand)
- [#613](https://github.com/deis/deis/pull/613) Merge pull request #613 from opdemand/localdev-docs (@opdemand)
- [#614](https://github.com/deis/deis/pull/614) Merge pull request #614 from opdemand/architecture-docs (@opdemand)
- [83e2b38](https://github.com/deis/deis/commit/83e2b3868583ab97899e9d51dcd2474a9ef69746) Removed old branch and service from .travis.yml. (@mboersma)
- [72ff95f](https://github.com/deis/deis/commit/72ff95fb2f02a14f774c00195ca64357f80725d8) make notes neon pink and warnings neon orange (@bacongobbler)
- [#616](https://github.com/deis/deis/pull/616) Merge pull request #616 from opdemand/update-travis-yml (@opdemand)
- [a961ea3](https://github.com/deis/deis/commit/a961ea3bc311cb54d5fd17e2beec0ab5b9d94e4d) Updated vagrant readme (@tscheepers)
- [3c5178a](https://github.com/deis/deis/commit/3c5178ae488494b9d9f1326a0157748521ffe526) Fixed grammer error (@tscheepers)
- [1b8af22](https://github.com/deis/deis/commit/1b8af22f92dd14463f27b09098855044cc377ee5) Set hostname of chefserver in provisioning (@tscheepers)
- [df32775](https://github.com/deis/deis/commit/df32775b4652818ca8089a8b71cf97714137c184) fix #542 (@bacongobbler)
- [#623](https://github.com/deis/deis/pull/623) Merge pull request #623 from tscheepers/patch-1 (@tscheepers)
- [#625](https://github.com/deis/deis/pull/625) Merge pull request #625 from opdemand/542-config-vars (@opdemand)
- [4393806](https://github.com/deis/deis/commit/43938065230e6f727d98dc22057333e19f4a52a0) move config test to before initial release (@bacongobbler)
- [55cc766](https://github.com/deis/deis/commit/55cc76661549eb154151de004c2f279562ec16dc) Updated the docs_requirements file, fixes #624. (@mboersma)
- [#626](https://github.com/deis/deis/pull/626) Merge pull request #626 from opdemand/624-fix-rest-docs (@opdemand)
- [#617](https://github.com/deis/deis/pull/617) Merge pull request #617 from opdemand/doc-message-colors (@opdemand)
- [ca72140](https://github.com/deis/deis/commit/ca72140ea89b2c615f8373becb27ce260f2ca781) Merge branch 'dsh-perm-fix' of https://github.com/tombh/deis into tombh-dsh-perm-fix (@mboersma)
- [6ab9df7](https://github.com/deis/deis/commit/6ab9df7628b05f6fa494b188ad1457f7ed22706d) Use sudo when invoking dshell/dsh command. (@mboersma)
- [6ba78b2](https://github.com/deis/deis/commit/6ba78b2b7f4e172b3fbf74cbb2fee8c60c70da97) raise BuildFormationError on build or delete (@bacongobbler)
- [abdf9bd](https://github.com/deis/deis/commit/abdf9bd2788bfa84491feb00f92aee8cd192c302) test new behaviour returns HTTP 400 (@bacongobbler)
- [#628](https://github.com/deis/deis/pull/628) Merge pull request #628 from opdemand/490-formation-without-credentials (@opdemand)
- [b3202fd](https://github.com/deis/deis/commit/b3202fd046ee0f8990fcc1ec657886c7090aaa3c) set min username length to 4 (@bacongobbler)
- [175fdf5](https://github.com/deis/deis/commit/175fdf530a3ca129ccb5b17adc98d82f6ed0560a) typo: must be >= 4 (@bacongobbler)
- [#594](https://github.com/deis/deis/pull/594) Merge pull request #594 from Springest/private_networking_support (@Springest)
- [49f0feb](https://github.com/deis/deis/commit/49f0feb2d21ee83fa1b95d965ad83d51e22c32c1) update dop to v0.1.6 (@bacongobbler)
- [#631](https://github.com/deis/deis/pull/631) Merge pull request #631 from opdemand/add-custom-user-auth (@opdemand)
- [8dde760](https://github.com/deis/deis/commit/8dde760afa7d7263e0d60ed6b0d273e77c6ed0ff) update registry submodule (@bacongobbler)
- [#633](https://github.com/deis/deis/pull/633) Merge pull request #633 from opdemand/update-registry (@opdemand)
- [3dcf63e](https://github.com/deis/deis/commit/3dcf63eda984590116ecb37699a2495a12d2b381) bump registry submodule (@bacongobbler)
- [742cd97](https://github.com/deis/deis/commit/742cd979b296400ea363e8d3f6c5fccdc8237dca) Update submodule SHAs. (@mboersma)
- [be65525](https://github.com/deis/deis/commit/be655252687afe649ab3753f5f74e1cbef3f4ae9) fix versioning typo for docs (@bacongobbler)
- [3c0c887](https://github.com/deis/deis/commit/3c0c887e8ca881992c347de621fd1ee6dde94949) added documentation on the default attributes (@bacongobbler)
- [b3881d0](https://github.com/deis/deis/commit/b3881d0b6a3528a30038de0646721d22a2391fa2) pin ruby to 2.0.0 (@bacongobbler)
- [#634](https://github.com/deis/deis/pull/634) Merge pull request #634 from opdemand/cookbook-docs (@opdemand)
- [02d2d1e](https://github.com/deis/deis/commit/02d2d1e2728b35d2f4e9758463467f7536a118ea) added rsyslog export documentation (@bacongobbler)
- [0ba71d0](https://github.com/deis/deis/commit/0ba71d09c71c364f30e2e3a8231e0aa2a3e88a63) Updated docs_requirements.txt, fixed warnings and api.admin docs. (@mboersma)
- [542f6d2](https://github.com/deis/deis/commit/542f6d2c8255d7f7d10b283a8ed93d5bc851c6ff) changed title to "Manage the Controller" (@bacongobbler)
- [#637](https://github.com/deis/deis/pull/637) Merge pull request #637 from opdemand/sphinx-cleanup (@opdemand)
- [#636](https://github.com/deis/deis/pull/636) Merge pull request #636 from opdemand/rsyslog-docs (@opdemand)
- [1de1ea2](https://github.com/deis/deis/commit/1de1ea222d9e28abb6ca65484bf7dea3126c43a5) Updated submodule SHAs. (@mboersma)
- [4163ef6](https://github.com/deis/deis/commit/4163ef6ac2099bf987370040210fe5963cd0538b) Updated CLI binary locations in README.rst. (@mboersma)

### v0.5.1 (2014/03/07 19:26 +00:00)
- [2df6f9a](https://github.com/deis/deis/commit/2df6f9ac1b3d5d415f69ac190747cbbc0d185577) Switch master to v0.5.1. (@mboersma)
- [58e007e](https://github.com/deis/deis/commit/58e007e969387c1e9b3d238b3f1b285732129398) Added django-guardian to docs_requirements.txt. (@mboersma)
- [6e1e2b4](https://github.com/deis/deis/commit/6e1e2b4b005b1e43fb6dde3e76fd631f59ad6d1f) Update Gemfile.lock (@mo-mughrabi)
- [#524](https://github.com/deis/deis/pull/524) Merge pull request #524 from mo-mughrabi/master (@mo-mughrabi)
- [cf79eab](https://github.com/deis/deis/commit/cf79eabb55517b925370b09c0b05300aa388e007) Added client binary support via PyInstaller, refs #472. (@mboersma)
- [6ffaeb6](https://github.com/deis/deis/commit/6ffaeb6a7727193cb4b0e6ed724dd8c2ba088274) Ease vagrant activation rule, now that port 8000 is default. (@mboersma)
- [048e389](https://github.com/deis/deis/commit/048e3899331f2bb2ce1b44a551a37a152f14a014) Removed default args from layers:update, fixes #525. (@mboersma)
- [#526](https://github.com/deis/deis/pull/526) Merge pull request #526 from opdemand/layers-update-bug-525 (@opdemand)
- [#527](https://github.com/deis/deis/pull/527) Merge pull request #527 from opdemand/pyinstaller-cli (@opdemand)
- [fc5b8b3](https://github.com/deis/deis/commit/fc5b8b3a4d558a8ebab1dccc887f86117cdf498f) Enabled Django ATOMIC_REQUESTS to map views to transactions. (@mboersma)
- [8214f74](https://github.com/deis/deis/commit/8214f748d1e1a341c41ddd3604d5ad01c5f54fcb) Clarify that gabrtv is BDFL for Deis decisions. (@mboersma)
- [#538](https://github.com/deis/deis/pull/538) Merge pull request #538 from opdemand/bdfl-notice (@opdemand)
- [b8bd170](https://github.com/deis/deis/commit/b8bd170054ee29a66bcb8969ac366a291f67cc07) change in name of prepare script. (@paulczar)
- [#539](https://github.com/deis/deis/pull/539) Merge pull request #539 from paulczar/patch-1 (@paulczar)
- [d82905d](https://github.com/deis/deis/commit/d82905d7035508cf8edf7d98f927d2573702e598) Updated chef-docker to 0.31.0. (@mboersma)
- [7f3685c](https://github.com/deis/deis/commit/7f3685c35e322446eacffaeb49a2b548fe2f44da) Updated deis/builder to v0.1.1. (@mboersma)
- [#537](https://github.com/deis/deis/pull/537) Merge pull request #537 from opdemand/474-atomic-requests (@opdemand)
- [6c902a6](https://github.com/deis/deis/commit/6c902a6cd53939b78036428ca6d2d65cfa3d9b9c) Revert ATOMIC_REQUESTS, fixes #541. (@mboersma)
- [9632606](https://github.com/deis/deis/commit/9632606b7c545eb882c95225e2aec5d8b56949ac) Adding port 8000 to the EC2 security group (@fagiani)
- [c1e23d5](https://github.com/deis/deis/commit/c1e23d554ccd5a202af879e5072cbbd525823fe1) `deis destroy` removes git remote only if it matches app, fixes #544. (@mboersma)
- [#545](https://github.com/deis/deis/pull/545) Merge pull request #545 from fagiani/master (@fagiani)
- [3ffc956](https://github.com/deis/deis/commit/3ffc95673c411ac4893eac8b0619b1ff7593f5c6) fix some inconstancies within the rackspace contrib instructions (@davidcollom)
- [#546](https://github.com/deis/deis/pull/546) Merge pull request #546 from davidcollom/rackspace_contrib (@davidcollom)
- [0abf453](https://github.com/deis/deis/commit/0abf4536034da03c2afaada83b34764c8aa6a646) Updated controller prep scripts to cache current Docker images. (@mboersma)
- [b8bf42b](https://github.com/deis/deis/commit/b8bf42b07f490780e4a003c037c45e6937301205) Added registry port 5000 to EC2 provisioning script, refs #550. (@mboersma)
- [be69e78](https://github.com/deis/deis/commit/be69e781f3b38465b4d419e5b89505f99adddb8c) Noted how to build CLI binaries for a release, fixes #472. (@mboersma)
- [#557](https://github.com/deis/deis/pull/557) Merge pull request #557 from opdemand/472-cli-binaries (@opdemand)
- [96b9074](https://github.com/deis/deis/commit/96b9074233a729a39cfc2918ec9ef638c9145c99) Updated deis-cookbook SHA. (@mboersma)
- [8912540](https://github.com/deis/deis/commit/89125404ea709d4c78cba84c96ee379ee5e4f28a) Removed 'converging formation' section, fixes #432. (@mboersma)
- [0cec11d](https://github.com/deis/deis/commit/0cec11d852e7e6980cab0ff20acc42d41a17c74f) Add knife params support to EC2 controller provisioner (@TiuTalk)
- [#559](https://github.com/deis/deis/pull/559) Merge pull request #559 from espnbr/feature/add-knife-params-support-to-ec2 (@espnbr)
- [67f9d78](https://github.com/deis/deis/commit/67f9d7884e7ed18879dca43a0f10ec0fffe1ccdb) Updated container submodules and cookbook SHA. (@mboersma)
- [f145a1f](https://github.com/deis/deis/commit/f145a1f3a323da700e3138ac56e4c9846fbd8594) Updated docker tags in image prep scripts. (@mboersma)
- [71779a2](https://github.com/deis/deis/commit/71779a25e94480589e485d75596a01db1d157372) Removed docker tags now that we're back to trusted builds. (@mboersma)
- [#566](https://github.com/deis/deis/pull/566) Merge pull request #566 from opdemand/remove-docker-tags (@opdemand)
- [7cbfc5d](https://github.com/deis/deis/commit/7cbfc5d60f98304041a601932ceca136134b04ff) Allow controller to run behind an SSL termination proxy (@jwilder)
- [85928c3](https://github.com/deis/deis/commit/85928c376384a63375b729ded0a39e37289f29a8) Updated EC2 AMIs for v0.5.1, fixes #567. (@mboersma)
- [#571](https://github.com/deis/deis/pull/571) Merge pull request #571 from opdemand/567-time-to-make-the-amis (@opdemand)
- [#570](https://github.com/deis/deis/pull/570) Merge pull request #570 from jwilder/jw-ec2-ssl (@jwilder)
- [91f1687](https://github.com/deis/deis/commit/91f1687b4c492a31f9988a296f64c2222fb9c1a2) Added comment about HTTPS proxy setting. (@mboersma)
- [f67bdee](https://github.com/deis/deis/commit/f67bdee07f18c36f1b4095ed58a6e626c7559c16) Updated Sphinx to v1.2.2. (@mboersma)
- [0171cc5](https://github.com/deis/deis/commit/0171cc5a04c98694903af24a6cdfd810dde2b026) Updated release procedure doc for upcoming v0.5.1. (@mboersma)
- [14bbd55](https://github.com/deis/deis/commit/14bbd5516e703ec3927bd2da1bdbc35689390201) removed provision-controller.sh (@bacongobbler)
- [#572](https://github.com/deis/deis/pull/572) Merge pull request #572 from opdemand/update-releases-doc (@opdemand)
- [#574](https://github.com/deis/deis/pull/574) Merge pull request #574 from opdemand/remove-provision-script (@opdemand)
- [5917999](https://github.com/deis/deis/commit/5917999e0552699f08d92469d29b2e66a2d62d40) Added missing DigitalOcean cloud regions. (@mboersma)
- [87096ec](https://github.com/deis/deis/commit/87096ecde1da0a53fba15e1e085b01f86001ac57) Always use :latest tag for Deis docker images. (@mboersma)
- [605e274](https://github.com/deis/deis/commit/605e274aaf1ab93194673123235a1cac4c67b076) Added :latest tag to slugrunner image pulls. (@mboersma)
- [#576](https://github.com/deis/deis/pull/576) Merge pull request #576 from opdemand/docker-tag-latest (@opdemand)
- [#575](https://github.com/deis/deis/pull/575) Merge pull request #575 from opdemand/digitalocean-regions (@opdemand)
- [31693ee](https://github.com/deis/deis/commit/31693ee3abaa249a324635fa3d74398905efe341) Updated submodules to match :latest tags. (@mboersma)
- [63954fc](https://github.com/deis/deis/commit/63954fcc843c95cf017193a0d75302c4fddc3126) Revert "removed provision-controller.sh" (@bacongobbler)
- [#579](https://github.com/deis/deis/pull/579) Merge pull request #579 from opdemand/remove-provision-script (@opdemand)
- [5b95c59](https://github.com/deis/deis/commit/5b95c5951f373926bf8852b4d16422c7a671c79e) Changes to Vagrant provisioner and environment to accomodate new (@tombh)
- [#528](https://github.com/deis/deis/pull/528) Merge pull request #528 from tombh/vagrant-provisioning (@tombh)
- [5083540](https://github.com/deis/deis/commit/508354092783f8eb85519d48a012e3dd9bb70263) fix merge comments from #528 (@bacongobbler)
- [#580](https://github.com/deis/deis/pull/580) Merge pull request #580 from opdemand/merge-fixes-from-528 (@opdemand)
- [bfaea45](https://github.com/deis/deis/commit/bfaea459c7a7644154e67e8fb77e33e0f2706940) change box_url to be the same as controller (@bacongobbler)
- [c9d6be0](https://github.com/deis/deis/commit/c9d6be0c470584dc19de4a158f190fba6f32fb16) change box name to deis-controller for consistency (@bacongobbler)
- [2814cce](https://github.com/deis/deis/commit/2814cce319fe54a408acc3beea0cd61919bd7c55) s/controller/server/g for consistency (@bacongobbler)
- [90dc2d8](https://github.com/deis/deis/commit/90dc2d81c34cf333409cdf9469426a1d3b9dc747) change default box_url. Fixes #581 (@bacongobbler)
- [#582](https://github.com/deis/deis/pull/582) Merge pull request #582 from opdemand/change-vagrantfile-template-url (@opdemand)
- [1a89091](https://github.com/deis/deis/commit/1a89091b05c78342b4c83ceb5d272060b8a14eb6) Updated EC2 AMIs from prep scripts. (@mboersma)
- [50193e8](https://github.com/deis/deis/commit/50193e892683c68d01b583e27d2271451a8cb1b0) Updated Chef references to 11.8.2. (@mboersma)
- [#584](https://github.com/deis/deis/pull/584) Merge pull request #584 from opdemand/chef-up (@opdemand)
- [#583](https://github.com/deis/deis/pull/583) Merge pull request #583 from opdemand/new-amis (@opdemand)
- [a2c194e](https://github.com/deis/deis/commit/a2c194e195dd9d8eba47b796d66b877a4c8d5235) fix buff-extensions downgrading (@bacongobbler)
- [#585](https://github.com/deis/deis/pull/585) Merge pull request #585 from opdemand/fix-buff-extensions (@opdemand)
- [aa25034](https://github.com/deis/deis/commit/aa25034cc707e06dd84e63417f46d60c55e8519a) Updated submodule SHAs. (@mboersma)
- [b1a8e14](https://github.com/deis/deis/commit/b1a8e14d1f08e542aa728de62694217df118c511) update DO readme to match correct script name (@davidcollom)
- [9cc0740](https://github.com/deis/deis/commit/9cc07409dd07c523ca7efb451cb819a9f042e5e0) be a little more explicit with the creation of snapshots (@davidcollom)
- [#588](https://github.com/deis/deis/pull/588) Merge pull request #588 from davidcollom/patch-1 (@davidcollom)
- [36a55b9](https://github.com/deis/deis/commit/36a55b91937812c95134887da524af61a1944dd5) pin knife-rackspace to 0.9.0 or higher (@bacongobbler)
- [43b170f](https://github.com/deis/deis/commit/43b170f234e6566f7807cbf41608bc9845095b23) Adds private networking for Digital Ocean (@Milodv)
- [66c93bc](https://github.com/deis/deis/commit/66c93bc0f0a672ad4a3264da8416050f227c3e5c) Updates dop version to support private networking and scrubbing (@Milodv)
- [ad784b7](https://github.com/deis/deis/commit/ad784b761a708b95e603f14c32533032d7bbc373) Updated deis-cookbook SHA to match master. (@mboersma)
- [9f7c761](https://github.com/deis/deis/commit/9f7c76131446221da974d96ed622a5e2d5d666d7) Fixed a typo in vagrant provisioning script. (@mboersma)
- [e7ad308](https://github.com/deis/deis/commit/e7ad308fe45fd0b2a11f395867c03e570aef5cfa) Use sudo to create symlink to dsh. #596 (@tombh)
- [723f985](https://github.com/deis/deis/commit/723f9852cc914cbade6b075d0ed3e2f527088e6c) added install script for avahi (@bacongobbler)
- [6c57c76](https://github.com/deis/deis/commit/6c57c7696032c51e826bf96396dd2d0cba453081) add avahi to other inline script instead (@bacongobbler)
- [85f8c1e](https://github.com/deis/deis/commit/85f8c1e5ed469a05427ac52ba9e25b467c3b7d95) make small improvements to the rackspace script (@bacongobbler)
- [#600](https://github.com/deis/deis/pull/600) Merge pull request #600 from opdemand/vagrant-avahi (@opdemand)
- [a54538a](https://github.com/deis/deis/commit/a54538a9667c6aa9d0598c61698bd33665311e61) remove redundant prepare script (@bacongobbler)
- [1350307](https://github.com/deis/deis/commit/135030730e3120ccf5fc911b032e7c5d6e3b3edd) use standardized format for naming nodes/layers (@bacongobbler)
- [b432280](https://github.com/deis/deis/commit/b432280b810e5752d4221ae790475f7be2335bca) append node id to end of name (@bacongobbler)
- [7de5caf](https://github.com/deis/deis/commit/7de5caf88ec03ecd17855871beec0ac51bcc3484) renamed prepare script (@bacongobbler)
- [9019f0d](https://github.com/deis/deis/commit/9019f0d526b0e537c93ada7b230707f9b80f0fb8) added nova client fo uploading SSH keys (@bacongobbler)
- [80e6a02](https://github.com/deis/deis/commit/80e6a028261cfd80d8147736fd2291e1945cd980) update Rackspace readme docs (@bacongobbler)
- [a5d8c05](https://github.com/deis/deis/commit/a5d8c0533094ab24315393851535dea7628953d1) bash script should source rackspacerc (@bacongobbler)
- [c6d4e43](https://github.com/deis/deis/commit/c6d4e4387e123343713c831f36806dee4d11dbfd) update prepare script reference (@bacongobbler)
- [eddea17](https://github.com/deis/deis/commit/eddea1702d35cd38436cc668e8b83c696da8cf11) make small improvements to the rackspace script (@bacongobbler)
- [cc805ee](https://github.com/deis/deis/commit/cc805ee8a592ebec331beb19ebbc42a7c33a9309) pin knife-rackspace to 0.9.0 or higher (@bacongobbler)
- [eed52c7](https://github.com/deis/deis/commit/eed52c70e34d4e8330856e93be0f1b18bee161f2) use standardized format for naming nodes/layers (@bacongobbler)
- [336d967](https://github.com/deis/deis/commit/336d967d3afb55c27a3df249a4cc111e3082cbfd) bash script should source rackspacerc (@bacongobbler)
- [196f996](https://github.com/deis/deis/commit/196f996bdb3bf534d52113a2c01372831c15a11d) added nova client fo uploading SSH keys (@bacongobbler)
- [f51eea3](https://github.com/deis/deis/commit/f51eea39f42f22c4c830f1287eb5c3e5a9947e1e) update prepare script reference (@bacongobbler)
- [145f7a2](https://github.com/deis/deis/commit/145f7a226aa80b3a2328b53692e8d3f8449cecba) append node id to end of name (@bacongobbler)
- [e31d9bd](https://github.com/deis/deis/commit/e31d9bd708d6fe7ba6632b4e83d7f85aeb956863) renamed prepare script (@bacongobbler)
- [7df9fcb](https://github.com/deis/deis/commit/7df9fcb5986cb40959b968aa215c7ffd939a0cf9) update Rackspace readme docs (@bacongobbler)
- [239ea2f](https://github.com/deis/deis/commit/239ea2fdbce2bce18564df266dd663e8e319ea19) remove redundant prepare script (@bacongobbler)
- [1e8395d](https://github.com/deis/deis/commit/1e8395de9511caa017f7e97d96221c85ab7dc24a) Merge branch 'fix-rackspace-provider' of github.com:opdemand/deis into fix-rackspace-provider (@bacongobbler)
- [3cea2ce](https://github.com/deis/deis/commit/3cea2cec985616b21cecfa934bc1095e7182aeea) Added description for releases:rollback. (@mboersma)
- [e3ab738](https://github.com/deis/deis/commit/e3ab73817dd4072fd28133e080229885a4d7df9a) Updated SHA to match deis-cookbook. (@mboersma)

### v0.5.0 (2014/02/13 04:18 +00:00)
- [d2e307d](https://github.com/deis/deis/commit/d2e307d4529d063a8dd2d6cb715701393e569f45) Switch master to v0.5.0. (@mboersma)
- [#503](https://github.com/deis/deis/pull/503) Merge pull request #503 from tombh/493-vagrant-detection (@tombh)
- [9885226](https://github.com/deis/deis/commit/9885226bfecdcbdcacd754597674fc4f4e1116af) rename gitignore to gitkeep
- [03bbaa1](https://github.com/deis/deis/commit/03bbaa1f69f99e99f6755f76b6c62456db45f7a5) Added containerize branch to Travis CI. (@mboersma)
- [#505](https://github.com/deis/deis/pull/505) Merge pull request #505 from opdemand/rename-gitignore (@opdemand)
- [0f83f6a](https://github.com/deis/deis/commit/0f83f6a995ff6df15819977ad87506c327983de2) Updated Django to v1.6.2. (@mboersma)
- [b85954d](https://github.com/deis/deis/commit/b85954d4ece060cba7d9b932d9255b92ae2ab3ff) Fixed build status badge to refer to master branch. (@mboersma)
- [98296b7](https://github.com/deis/deis/commit/98296b72540c5c0a82ee40b0aa9bd7e005fb833f) add dynamic confd settings lookup (with local_settings override) (@gabrtv)
- [3dfe59c](https://github.com/deis/deis/commit/3dfe59cbe6b575810fc715d29ac79618cecef724) install etcd gem prior to chef run, cleanup cruft (@gabrtv)
- [b8a5518](https://github.com/deis/deis/commit/b8a5518efeb693a686b9fa7b0e572fe981159d22) shift ssh:// to https:// for read-only access (@gabrtv)
- [37a3fe6](https://github.com/deis/deis/commit/37a3fe6928f38f5f90a4eef24e6d8073f4ed592d) bump docker images (@gabrtv)
- [2acdd2a](https://github.com/deis/deis/commit/2acdd2aaeec216c2fbf39befb2ff6caad04a475a) new hooks API for push/build, deprecate Build.push() (@gabrtv)
- [323af3b](https://github.com/deis/deis/commit/323af3b355339ddff2e9a49fc92c46ae16e192e1) update image SHAs (@gabrtv)
- [17aaac7](https://github.com/deis/deis/commit/17aaac7358dcd684823dc110a4f0e24e474682ac) remove builder debug (@gabrtv)
- [15a71a8](https://github.com/deis/deis/commit/15a71a803291cc97383116be71568fe72343311c) handle builder target that contains a port (@gabrtv)
- [5af0968](https://github.com/deis/deis/commit/5af0968ce71fea52d9539a08eeb903d0c83ec86e) stop using sudo in chef setup (@gabrtv)
- [b2858fe](https://github.com/deis/deis/commit/b2858fea970bb9bc9dda6694bcf08711f44527c8) switch from rabbitmq to redis, install gunicorn app for run_gunicorn support (@gabrtv)
- [1b1875b](https://github.com/deis/deis/commit/1b1875b8effc62dd4fc9223a849f3774814b47a5) new redis 2.8 dependency (@gabrtv)
- [fd9bbcf](https://github.com/deis/deis/commit/fd9bbcfbcb3d051cb865cc5d82f8d0230b8da257) ignore import errors on local_settings (@gabrtv)
- [b98e4d8](https://github.com/deis/deis/commit/b98e4d8b988ed22b8ef4c07cf645adee4b43d4bb) new vagrantfile for container-based deployment (@gabrtv)
- [0331f51](https://github.com/deis/deis/commit/0331f51677ab089b36e09139b849181b6bed0e51) initial pass at architecture docs (@gabrtv)
- [81c9ca4](https://github.com/deis/deis/commit/81c9ca4bb85b0774ce739d9452a66f259763e1ae) add docker images as submodules (@gabrtv)
- [f5312c5](https://github.com/deis/deis/commit/f5312c58ef7ae197955fe4d49d8cbc0eb7b57bc1) switch git remote to ssh:// syntax with 2222/tcp (@gabrtv)
- [d42cee2](https://github.com/deis/deis/commit/d42cee226058568dbd982f2fbc11037b21e76d44) change vagrant base box to vanilla ubuntu 12.04 w/ 3.8 kernel (@gabrtv)
- [b467998](https://github.com/deis/deis/commit/b4679982b614e22efde21734bd520587836e5021) fix vagrant provider in containerized env by using ip_addr instead of avahi hostnames (@gabrtv)
- [c1f4dc5](https://github.com/deis/deis/commit/c1f4dc524c4596d4bd91facc12fad4c5abce295b) pre-install required etcd gem before knife bootstrap using template-file option (@gabrtv)
- [b6c4451](https://github.com/deis/deis/commit/b6c4451ca21824b6b03d1ed38921eefa32188171) change run to use build.image instead of a slug mount (@gabrtv)
- [9bdba6c](https://github.com/deis/deis/commit/9bdba6ccb8533e6bfe230385fe09904f56f21f46) deprecate converge_controller (@gabrtv)
- [bd31773](https://github.com/deis/deis/commit/bd31773dc2faf7f2e35b16aa6de31c1223419bcf) update container SHAs (@gabrtv)
- [5a1eb0b](https://github.com/deis/deis/commit/5a1eb0beccc9dfb3884a48b1f2920f41596302fe) add rsyslog logger component (@gabrtv)
- [ae37b11](https://github.com/deis/deis/commit/ae37b11105f62d24b8ea435ded752ddd095e0df8) add vagrant authorized-keys setup script (@gabrtv)
- [b33775f](https://github.com/deis/deis/commit/b33775f3d87f937f14ba8cddd2ba8ff4f725c544) switch to containerized cookbooks branch (@gabrtv)
- [af2785a](https://github.com/deis/deis/commit/af2785af9da3ae20f5a4a79b99f03961d45d5b0b) remove knife template for etcd gem at load-time (@gabrtv)
- [bd89618](https://github.com/deis/deis/commit/bd8961843c8ecdfd88b3ed24051d5408df18a731) add dev mode/source flags to Vagrantfile (@gabrtv)
- [eca3d9f](https://github.com/deis/deis/commit/eca3d9fcde4db6ce8adc4c879811c0ebbf0b2bf8) separate controller/node ec2 ami prep (@gabrtv)
- [bd2ba0a](https://github.com/deis/deis/commit/bd2ba0ad9ebdb450a98c35e5d399e8f985f8f2ab) publish SSH Keys to etcd instead of a chef data bag (@gabrtv)
- [405245b](https://github.com/deis/deis/commit/405245b749dfe5033c5f886c194cb4044b6c949b) deprecate deis-users databag (@gabrtv)
- [4edac55](https://github.com/deis/deis/commit/4edac55531433363ec0a09e498e31f8d6a52442b) update docker image SHAs (@gabrtv)
- [56c28e1](https://github.com/deis/deis/commit/56c28e18461311208e9cb5eca135c17f89a065ff) Updated DigitalOcean scripts for containerize support. (@mboersma)
- [08fd115](https://github.com/deis/deis/commit/08fd115562aa83fb3a05852c956c77af0adc4089) Fixed flake8 errors. (@mboersma)
- [cf4f174](https://github.com/deis/deis/commit/cf4f174dabc8a632512cc00f1a8209ebb93f1b7b) Updated Rackspace scripts for containerize support. (@mboersma)
- [c82b8fb](https://github.com/deis/deis/commit/c82b8fb7908f8edb53cb04be55153976036df336) update build and server images (@gabrtv)
- [3913498](https://github.com/deis/deis/commit/391349860d2de18ad44b1a1c8216f8e776c1b3fc) return app info on build hook (@gabrtv)
- [c7251a7](https://github.com/deis/deis/commit/c7251a73754038d018e1ea94d5db6101c112cad0) spaces > tabs (@gabrtv)
- [ae23502](https://github.com/deis/deis/commit/ae235021705434571893f91eff8a217ec63cc218) fix `deis run` target image (@gabrtv)
- [a3c2124](https://github.com/deis/deis/commit/a3c212426b543e19dc4a2c81bc32686585cd8a0a) update registry and server images (@gabrtv)
- [1b7e1f3](https://github.com/deis/deis/commit/1b7e1f3c7d25095ce20d91a8521db2ea570c2b44) publish every release as a tagged docker image (ex: gabrtv/myapp:v23) (@gabrtv)
- [8afa178](https://github.com/deis/deis/commit/8afa178973173fff468ed05755d6e95dd163b3fc) update server image w/ missing logs directory (@gabrtv)
- [daafe17](https://github.com/deis/deis/commit/daafe175e47594044567a4cf4357ce853cd1aaca) upgrade to chef-docker 0.28.0 for new :pull behavior (@gabrtv)
- [27f483c](https://github.com/deis/deis/commit/27f483ce13fb080ac27ec30b1725616ce9635e57) Added containerize branch to Travis CI. (@mboersma)
- [b40aa2a](https://github.com/deis/deis/commit/b40aa2a5a9a973d16b5ce486f8057c9c969eed0e) remove dupe release, fix rollback publishing (@gabrtv)
- [27c1efe](https://github.com/deis/deis/commit/27c1efe698a3a42c190b29e1f30fb7389f28f613) include sha on every git push event (@gabrtv)
- [1427d75](https://github.com/deis/deis/commit/1427d755d77ec666d08e94ec12a17c8e3483aaa3) Updated SHA for deis-cookbook containerize branch. (@mboersma)
- [89e649e](https://github.com/deis/deis/commit/89e649edb9185362957418d67b4bae950e72d083) Updated Django to v1.6.2. (@mboersma)
- [0fb67df](https://github.com/deis/deis/commit/0fb67df19df13cf50687d1bdedd68289c54ee911) Updated requirements files and deis/server SHA. (@mboersma)
- [2a5439e](https://github.com/deis/deis/commit/2a5439e40c01b4fc7327877ab053c1fc38808971) Added static file serving WSGI app, fixes #507. (@mboersma)
- [2693349](https://github.com/deis/deis/commit/2693349d3181cdbbbf6ff0c843cd6599c5dd86c9) rename submodules to match container names (@gabrtv)
- [31d455f](https://github.com/deis/deis/commit/31d455f831ae1aaf58d4bb9e1dca4080bb913509) update server/worker SHAs (@gabrtv)
- [ef0dcff](https://github.com/deis/deis/commit/ef0dcff5db61804fafe6ebe12ca996937150e1a6) Removed references to gitosis. (@mboersma)
- [008042b](https://github.com/deis/deis/commit/008042b728eb70294f1c2f8c881a5b827d4cfbcf) added vagrant plugin commands
- [29733f7](https://github.com/deis/deis/commit/29733f7e49d029be83b1917bb421237fcd7ed5aa) remove tab in vagrantfile
- [0f2b688](https://github.com/deis/deis/commit/0f2b68863e1cbaff1b3278dc6f363b3aaea49150) Removed references to obsolete deis-users data bag. (@mboersma)
- [b5a55c5](https://github.com/deis/deis/commit/b5a55c54b305150da58719a3a36787f131193aa5) new images (@gabrtv)
- [4e94fce](https://github.com/deis/deis/commit/4e94fce119fa02e76f2db3b8989c46d3c7280a39) Updated SHAs for submodules and cookbook. (@mboersma)
- [b5afa5a](https://github.com/deis/deis/commit/b5afa5a230e2baf38013bacbc1ffce39d00e1edb) Dockerfile support no longer 'coming soon'. (@mboersma)
- [c5ea9b2](https://github.com/deis/deis/commit/c5ea9b28200b15d751cdb2fd26f632e6f9cdadbe) Updated SHAs to latest deis projects. (@mboersma)
- [4fbe908](https://github.com/deis/deis/commit/4fbe908a46d26cc27bafd7a41194743193577cb7) assert flavor count is > 0 (@gabrtv)
- [9f0ad15](https://github.com/deis/deis/commit/9f0ad156cbaf26d2ebaea5fa476ef7d72847f0b3) remove deprecated key CM test (@gabrtv)
- [99b236a](https://github.com/deis/deis/commit/99b236ac606abe653adc105ce25495a5b979aa7f) defaut to sqlite3 if not defined (@gabrtv)
- [2c923fd](https://github.com/deis/deis/commit/2c923fd9faecbc392994cc621459833c088f52ac) move PROVIDER_MODULES into component images (@gabrtv)
- [c942eec](https://github.com/deis/deis/commit/c942eec4e7500f1b76fe2cdb9c7d8af678fe8d24) fix hook tests and import into test runner (@gabrtv)
- [b14246b](https://github.com/deis/deis/commit/b14246ba34f72f321faae939b660a00329421dfc) rename default database to deis.db (@gabrtv)
- [795bf9d](https://github.com/deis/deis/commit/795bf9d753378b771dbd92b846d4ec34a911fb28) add test for ssh key fingerprinting (@gabrtv)
- [129a381](https://github.com/deis/deis/commit/129a381fd9ec01780e77ae99c59abe070275da29) exclude osx system libs and docker registry integration from coverage (@gabrtv)
- [859646f](https://github.com/deis/deis/commit/859646f9bff344fedf78892b0ed301f2f6fbba3f) flake8 fixes (@gabrtv)
- [78c7696](https://github.com/deis/deis/commit/78c76967ddc20fb0cfec61e0833699957c028c89) update server and worker w/ v0.5.0 tag (@gabrtv)
- [#522](https://github.com/deis/deis/pull/522) Merge pull request #522 from opdemand/containerize (@opdemand)
- [1b7033c](https://github.com/deis/deis/commit/1b7033c2362b6322f4e7ff4f1f895d16917fc098) Fixed merge remnant in Berksfile.lock. (@mboersma)
- [bbe2ea9](https://github.com/deis/deis/commit/bbe2ea930173fc3b1d231cb7089dc264c4f5dd9a) Updated deis-cookbook SHA. (@mboersma)
- [58c0429](https://github.com/deis/deis/commit/58c042955d8fa4a8a6688b80792b3e2d3647bc15) re-target Berksfile to deis-cookbook master, update other cookbooks (@gabrtv)
- [9cacbac](https://github.com/deis/deis/commit/9cacbac5b3bccf83a9c7b42c7dd7a6bfe6b93c91) Updated image prep scripts. (@mboersma)

### v0.4.1 (2014/02/04 16:42 +00:00)
- [44a73c3](https://github.com/deis/deis/commit/44a73c38a7ab40ffb883fcdf0bd8d804a2dae640) Fixed required libs in setup.py to use newer install_requires. (@mboersma)
- [36b6a9c](https://github.com/deis/deis/commit/36b6a9c3bb691693c7a6069b9569cb5ae8b32300) Switch master to v0.4.1. (@mboersma)
- [bba5f3d](https://github.com/deis/deis/commit/bba5f3d1111a4e64f40dca4d075118dfb0140bae) Updated djangorestframework to 2.3.12 security fix. (@mboersma)
- [0317227](https://github.com/deis/deis/commit/031722771ea810f80e78dc247b6823f5c3eee3f0) Updated Docker to v0.7.6. (@mboersma)
- [#463](https://github.com/deis/deis/pull/463) Merge pull request #463 from opdemand/docker-0.7.6 (@opdemand)
- [86e8268](https://github.com/deis/deis/commit/86e826802eb26877b8d4387fe379672759801315) Show the list of config variables sorted by name (@nathansamson)
- [6572323](https://github.com/deis/deis/commit/657232350be6a2be45acac369e5da66cf6be69e0) Tabularize config output (@nathansamson)
- [7419a54](https://github.com/deis/deis/commit/7419a54527980b1466782d4095277650b9f7e815) Add a --oneline option to config:list (@nathansamson)
- [#469](https://github.com/deis/deis/pull/469) Merge pull request #469 from nathansamson/nathan/client-config-improvements (@nathansamson)
- [ccbf33d](https://github.com/deis/deis/commit/ccbf33dc9c94b934bc0ef9cccba1fcb873e59362) significant typo ("no" -> "yes") (@jfw)
- [0ba9ed1](https://github.com/deis/deis/commit/0ba9ed1bdb726cdc8ca03c3ff6c2f28b4ca61a18) Allow the provision-controller.sh script to run in a directory with spaces (@nathansamson)
- [#471](https://github.com/deis/deis/pull/471) Merge pull request #471 from nathansamson/nathan/vagrant-space-fix (@nathansamson)
- [#470](https://github.com/deis/deis/pull/470) Merge pull request #470 from jfw/patch-1 (@jfw)
- [fdb1b30](https://github.com/deis/deis/commit/fdb1b30c0e3b8ed29a10cb5a4df4b1ce54a7e3c3) Fix another space-in-path issue in vagrant provisioner (@nathansamson)
- [410ebac](https://github.com/deis/deis/commit/410ebac516277b0c6a82967d77470623088833a6) Fix spaces-in-path issues for DigitalOcean contrib scripts. (@nathansamson)
- [8201066](https://github.com/deis/deis/commit/8201066250ec05a568ffee229d1cce3c745508d3) Fix EC2 provision script (space in path issue) (@nathansamson)
- [21ac978](https://github.com/deis/deis/commit/21ac978bcdab66681ef0882fc28e4cc9018a1e33) Fix space issue for rackspace controller (@nathansamson)
- [db9404e](https://github.com/deis/deis/commit/db9404ea3e6016f5ef79de380c43ef4e491628e1) Allow vagrant (host sytem) to run in a directory with spaces (@nathansamson)
- [5c0ddb6](https://github.com/deis/deis/commit/5c0ddb6e0690cf579fd3c4a9b1df40a6a0798249) For dev convenience, added a DB-reset script and some basic fixtures. Also attempt to deal with git-ignoring 'static' (@tombh)
- [#473](https://github.com/deis/deis/pull/473) Merge pull request #473 from nathansamson/nathan/vagrant-space-fix2 (@nathansamson)
- [#479](https://github.com/deis/deis/pull/479) Merge pull request #479 from nathansamson/vagrant-fixes (@nathansamson)
- [#483](https://github.com/deis/deis/pull/483) Merge pull request #483 from tombh/db-reset-script (@tombh)
- [bd0b0dc](https://github.com/deis/deis/commit/bd0b0dc1ac3812a8a5162e4c3aaa97ceba716054) DigitalOcean: Fix for #477 (@nathansamson)
- [#487](https://github.com/deis/deis/pull/487) Merge pull request #487 from nathansamson/nathan/477 (@nathansamson)
- [c2951f9](https://github.com/deis/deis/commit/c2951f91d19d65b415777754ab3a681b70b4c411) Showing containers for the correct app when running scale --app=X. Fixes #481 (@nathansamson)
- [bbda0a6](https://github.com/deis/deis/commit/bbda0a69999673bc8789e08abc2f36786dc119ef) Allow users to patsh to enters SSH keys in the CLI client (@nathansamson)
- [0dc199f](https://github.com/deis/deis/commit/0dc199fcf16f658b3a9b4810f85ec1f30ee9803f) Correctly honor args for deis config comand (@nathansamson)
- [7af2427](https://github.com/deis/deis/commit/7af24272a4ebaa6b4b89141a1bb86c681c2e923c) Correctly honor args for deis perms comand (@nathansamson)
- [9d23b0f](https://github.com/deis/deis/commit/9d23b0f7969b3735c29adb335c2ca0351ab0e27c) Change the arg parsing for fallback :list commands according to a found example in the code (@nathansamson)
- [7f3b95f](https://github.com/deis/deis/commit/7f3b95f836bd1a3ca7246bf848c9b1d31ff512f4) change `deis run` bind-mount to read-only to prevent disruptive modification (@gabrtv)
- [#489](https://github.com/deis/deis/pull/489) Merge pull request #489 from opdemand/fix-deis-run (@opdemand)
- [93f8210](https://github.com/deis/deis/commit/93f82102af8d85ecc1e443b8fb9cef3f1bb77c59) Trick flake into decreasing keys_add CC (@nathansamson)
- [#488](https://github.com/deis/deis/pull/488) Merge pull request #488 from nathansamson/nathan/cli-improvements (@nathansamson)
- [0e768b8](https://github.com/deis/deis/commit/0e768b84682218e87a8c8f0cf39f11e554e050ee) clean up /etc/chef, install inotify-tools (@paulczar)
- [1d1c661](https://github.com/deis/deis/commit/1d1c66105441097804f32f43ba3f63d5ec0018b5) Clean up /etc/chef and install inotify-tools across all providers. (@mboersma)
- [b7751ab](https://github.com/deis/deis/commit/b7751ab212d1abf93bf245efe688772a027fa4a7) Removed deprecated pip --use-mirrors flag. (@mboersma)
- [a355a0e](https://github.com/deis/deis/commit/a355a0e487ce95ac18a81bf63fdffc8af336bd24) Updated celery, requests, Sphinx. (@mboersma)
- [d601fd8](https://github.com/deis/deis/commit/d601fd8bf2b9cbff3923efd290ea2b5fca34a008) Updated requests version in setup.py to 2.2.1. (@mboersma)
- [169627f](https://github.com/deis/deis/commit/169627f63563a595f7f062da12fa4705696927c9) Add formation, provider and flavours to fixtures for more complete seeding (@tombh)
- [#491](https://github.com/deis/deis/pull/491) Merge pull request #491 from opdemand/pypi-updates (@opdemand)
- [45aff6c](https://github.com/deis/deis/commit/45aff6ca665a11aacbbfaa1c2826a5cccd5ab0e3) Updated Berksfile.lock to latest deis-cookbook (@mboersma)
- [cf6f74f](https://github.com/deis/deis/commit/cf6f74f9efe613ba2454e52c24756851d987fd7a) Only take into account the first equal sign when setting config vars (@nathansamson)
- [42d1d0b](https://github.com/deis/deis/commit/42d1d0b11a6205dfe51208ca07922c4b29fd36e8) Updated knife-digital_ocean to v0.4.0. (@mboersma)
- [#500](https://github.com/deis/deis/pull/500) Merge pull request #500 from nathansamson/nathan/equalsinvars (@nathansamson)
- [f804ed8](https://github.com/deis/deis/commit/f804ed8323ba99a2d9d48f648089ff2c42c2ea3b) Revert bind-mounting slugs read-only. (@mboersma)
- [#502](https://github.com/deis/deis/pull/502) Merge pull request #502 from opdemand/revert-pr-489 (@opdemand)
- [#501](https://github.com/deis/deis/pull/501) Merge pull request #501 from opdemand/knife-do-update (@opdemand)
- [#494](https://github.com/deis/deis/pull/494) Merge pull request #494 from tombh/db-reset-script (@tombh)
- [4e79646](https://github.com/deis/deis/commit/4e796464c15acaddd09322ee91f0401eb350b1a3) Added detailed steps to add deis-controller to admin group (@shredder12)
- [a2ad8a0](https://github.com/deis/deis/commit/a2ad8a0651a6b2e84c1ce89350e0149dc83a0f24) Added link to Chef admins edit, admonition section. (@mboersma)
- [2d10965](https://github.com/deis/deis/commit/2d1096508fbf1f63fcd49c210e8beaadf4ee1fdf) Replace method for checking if client is in a vagrant setup. Use (@tombh)

### v0.4.0 (2014/01/14 19:19 +00:00)
- [a5731be](https://github.com/deis/deis/commit/a5731be3dcd288c854c253294a0d2e3d9a2fc767) Switch master to v0.4.0. (@mboersma)
- [e1d60b2](https://github.com/deis/deis/commit/e1d60b2e392365cf34f1d28d2fb97afb7e245c79) Restored pypip.in download badge, added license badge. (@mboersma)
- [2b4a8db](https://github.com/deis/deis/commit/2b4a8db663d053d4c4819a9f17ad84daa08df903) Added reminder to rebuild all published readthedocs versions on release. (@mboersma)
- [f22a6e5](https://github.com/deis/deis/commit/f22a6e573227f4cb0bf17801f2968471e59fe24a) Updated DigitalOcean snapshot script in line with other providers. (@mboersma)
- [ac35cd8](https://github.com/deis/deis/commit/ac35cd859ed55c6369f5df5836f32241ee40b12a) Updated boto to 2.21.2 (@mboersma)
- [75674ac](https://github.com/deis/deis/commit/75674acdd6e7e5c2c2d88ae0765aed4d4ae6bc04) Simplified IsAnonymous test. (@mboersma)
- [c1e60ee](https://github.com/deis/deis/commit/c1e60ee1008602ae50235f58eab7a10ae810d69b) Added docstrings for the decorated api.tasks module. (@mboersma)
- [9f9f861](https://github.com/deis/deis/commit/9f9f86122256225db8f29db09a69452174a0cf01) Improved CSS styling for [source] / [docs] links in documentation. (@mboersma)
- [01928c4](https://github.com/deis/deis/commit/01928c48462edb52b6aeb8cc4df20f649916f7cc) Cleaned up a few docstrings and removed dead South introspections. (@mboersma)
- [8d8daef](https://github.com/deis/deis/commit/8d8daef0ad890a21b3c6e42ee07d8fc0e9d9ea5d) Fixed #435 -- document how to use custom buildpacks in the FAQ. (@mboersma)
- [#437](https://github.com/deis/deis/pull/437) Merge pull request #437 from opdemand/more-docstrings (@opdemand)
- [ee796fe](https://github.com/deis/deis/commit/ee796fe723c4c444beca21dbf9cadccddc436e16) Updated Docker to v0.7.3. (@mboersma)
- [4646a95](https://github.com/deis/deis/commit/4646a95450a4686f12dd65b364423076c6caf1e4) Fixed typo (@joelvh)
- [#443](https://github.com/deis/deis/pull/443) Merge pull request #443 from opdemand/docker-0.7.3 (@opdemand)
- [#442](https://github.com/deis/deis/pull/442) Merge pull request #442 from joelvh/patch-1 (@joelvh)
- [#410](https://github.com/deis/deis/pull/410) Merge pull request #410 from tombh/376-move-slugbuilder-hook (@tombh)
- [a3aa6d4](https://github.com/deis/deis/commit/a3aa6d4a0787ba255ebbbc65e6c68f47c993954b) Swap steps 4 and 5 in vagrant README. Add step to use Makefile when installing deis client. (@tombh)
- [#445](https://github.com/deis/deis/pull/445) Merge pull request #445 from tombh/improve-vagrant-docs-for-niccolox (@tombh)
- [5d0deae](https://github.com/deis/deis/commit/5d0deae9f26992ef7956147a58322ddc2941de2c) Fixed #449 -- `deis releases:info` honors specified version. (@mboersma)
- [#450](https://github.com/deis/deis/pull/450) Merge pull request #450 from opdemand/449-get-release-version (@opdemand)
- [f7d887d](https://github.com/deis/deis/commit/f7d887dda25e003908ab9377ef74372a472d7e09) Added support for `deis release:rollback`, fixes #382. (@mboersma)
- [#451](https://github.com/deis/deis/pull/451) Merge pull request #451 from opdemand/382-release-rollback (@opdemand)
- [4586f61](https://github.com/deis/deis/commit/4586f61d7735a453678ec01833c3499c2eea36d1) Remove Authors from README (@gabrtv)
- [18f1810](https://github.com/deis/deis/commit/18f18108de066f45f23765808609710d3bc607d6) push implementation and tests for external builder module (@gabrtv)
- [f73b941](https://github.com/deis/deis/commit/f73b941ddef2678ca1c38a24657c1252cc7e55d4) remove providers from coveragerc, while continuing to import them (@gabrtv)
- [#453](https://github.com/deis/deis/pull/453) Merge pull request #453 from opdemand/docker-0.7.5 (@opdemand)
- [#452](https://github.com/deis/deis/pull/452) Merge pull request #452 from opdemand/external-push (@opdemand)
- [2f039e7](https://github.com/deis/deis/commit/2f039e7f4c7a20caef6257b328293e3aa8ea6e87) Fixed #383 -- Added summaries to `deis releases:list` (@mboersma)
- [263df70](https://github.com/deis/deis/commit/263df706b1a2e5ad252e0ae7ce58e9e8449bc6cb) Get returns current Config, not latest Config. Fixes #455. (@mboersma)
- [bf4de12](https://github.com/deis/deis/commit/bf4de12937044829641d0054d402ef1de9af0181) Reduce complexity in deis.py's main() by refactoring command dispatch lines (@tombh)
- [#456](https://github.com/deis/deis/pull/456) Merge pull request #456 from opdemand/455-get-config (@opdemand)
- [aa6f4c2](https://github.com/deis/deis/commit/aa6f4c24d5ec9b733f2aaf8f2b61f1e9a93f0f77) Added syslog events for app lifecycle. (@mboersma)
- [#446](https://github.com/deis/deis/pull/446) Merge pull request #446 from opdemand/383-releases-list (@opdemand)
- [#444](https://github.com/deis/deis/pull/444) Merge pull request #444 from tombh/client-feedback-for-unresponsive-controller (@tombh)
- [#457](https://github.com/deis/deis/pull/457) Merge pull request #457 from opdemand/394-app-lifecycle-events (@opdemand)
- [fe5a2c3](https://github.com/deis/deis/commit/fe5a2c3a28cce51f2c72db469d9dbb3ae205145d) Updated boto, paramiko, psycopg2, requests. (@mboersma)
- [#460](https://github.com/deis/deis/pull/460) Merge pull request #460 from opdemand/package-updates (@opdemand)
- [ad22b6c](https://github.com/deis/deis/commit/ad22b6ce8af5031db6a4bbf006d745cfadcd879e) Add the step of prepare a new image to the Rackspace contrib README (@stackedsax)
- [8d21adc](https://github.com/deis/deis/commit/8d21adc3030ed788dfb4b8e83c4efaac5cb0ac0d) upgrade the Rackspace prepare, provision, and provider scripts to use newer Rackspace performance flavors (@stackedsax)
- [6a32c8e](https://github.com/deis/deis/commit/6a32c8e078159f1ae07287245ef9bc9286e359da) Make the instructions simpler and clearer on how to run the prepare-rackspace-image.sh script.  Run from curl | bash (@stackedsax)
- [0e5312a](https://github.com/deis/deis/commit/0e5312a310eebe7d694ba55176cac7ebe6463995) Add a note to direct people to the old OpsCode control panel to add deis-controller to the admins group. (@stackedsax)
- [c1721ac](https://github.com/deis/deis/commit/c1721ac6425121d76894a7795056c9c5c54cf04f) Fixed test_auth after Rackspace changes removed one default flavor. (@mboersma)
- [4ca1c43](https://github.com/deis/deis/commit/4ca1c43e25e4830f207d7439adfd50a348973f56) Updated EC2 AMIs for v0.4.0, fixes #458 (@mboersma)

### v0.3.1 (2013/12/31 16:09 +00:00)
- [82a4799](https://github.com/deis/deis/commit/82a4799f21561301c8058a36cd7fc09514f9b245) Switch master to v0.3.1 (@mboersma)
- [78efe2c](https://github.com/deis/deis/commit/78efe2c210435615288a1a8a0d3c05aa9f80c913) Make $node_name unique to avoid deleting wrong droplet on failed provision (@tombh)
- [c91da1b](https://github.com/deis/deis/commit/c91da1b1bcf63ee44f6d902d488137fe7473210b) Removed .gitignore entry that prevented adding static assets. (@mboersma)
- [12ab065](https://github.com/deis/deis/commit/12ab0650c6aeb26dc7940f96c9f52727e829a344) Updated docs styling to avoid theme breakage on mobile devices (@bengrunfeld)
- [57d7a50](https://github.com/deis/deis/commit/57d7a50db257e0e890aed0d36d5375e23b9ede47) Added navigation elements for docset versions. Refs #256. (@mboersma)
- [#393](https://github.com/deis/deis/pull/393) Merge pull request #393 from opdemand/rtfd (@opdemand)
- [0bfba55](https://github.com/deis/deis/commit/0bfba55d9252968cb47ddd405d0cf54693db807c) Updated to Django 1.6.1. (@mboersma)
- [#395](https://github.com/deis/deis/pull/395) Merge pull request #395 from opdemand/django-1.6.1 (@opdemand)
- [3a627ed](https://github.com/deis/deis/commit/3a627ed8b4eb54c7567c7cf7b748017edbf0cdd0) Added Bing Webmaster Tools validation token (@bengrunfeld)
- [#392](https://github.com/deis/deis/pull/392) Merge pull request #392 from tombh/digital-ocean-provision-error-tolerance (@tombh)
- [6bd1fc1](https://github.com/deis/deis/commit/6bd1fc1c0006b3d96fa880ef05242d9cf217b170) Tweak tr command to work on OSX. Chef object deletion should use $node_name. opdemand/deis#396 (@tombh)
- [#401](https://github.com/deis/deis/pull/401) Merge pull request #401 from tombh/digital-ocean-provision-error-tolerance (@tombh)
- [c6d986a](https://github.com/deis/deis/commit/c6d986acdacf7ed527323f6fe7f4454fcc47911b) Moved slugbuilder from cookbook to deis. opdemand/deis#376 (@tombh)
- [f581767](https://github.com/deis/deis/commit/f5817672b2e64d01c25a291e3e1c5ef0ec690b3a) When checking the error type from a failed vagrant desutrcution, make comparison case insensitve - for BASH and ZSH support. opdemand/deis#346 (@tombh)
- [953e329](https://github.com/deis/deis/commit/953e3290d9fc65de330418b50bd67aa9afb9305b) Change Installation link to Operations Guide (@gabrtv)
- [#407](https://github.com/deis/deis/pull/407) Merge pull request #407 from tombh/346-node-del-check-case-insensitive (@tombh)
- [a3ac632](https://github.com/deis/deis/commit/a3ac6327ca7c24560128267dcd344164b3a78704) Updated several python packages. (@mboersma)
- [8e17497](https://github.com/deis/deis/commit/8e1749770025bf5ecd94481f68d80251249dd15d) Remove check for set env vars os.environ() does that anyway (@tombh)
- [#411](https://github.com/deis/deis/pull/411) Merge pull request #411 from opdemand/pypi-updates (@opdemand)
- [a2b446c](https://github.com/deis/deis/commit/a2b446cb4273c3fee8baf73f9566de726b0cdedf) Be specific about required package versions in client setup.py. (@mboersma)
- [0896b9b](https://github.com/deis/deis/commit/0896b9b562180368f143d767002a61ca064bf653) Fixed flake8 errors (@tombh)
- [a8e19d3](https://github.com/deis/deis/commit/a8e19d3e05f680bc937f39ac029c8b9901b5b76c) Updated README.rst for Deis CLI, fixes #409. (@mboersma)
- [#417](https://github.com/deis/deis/pull/417) Merge pull request #417 from opdemand/fix-409-client-readme (@opdemand)
- [#414](https://github.com/deis/deis/pull/414) Merge pull request #414 from opdemand/fix-413-requests-version (@opdemand)
- [676aa57](https://github.com/deis/deis/commit/676aa571516d5cad9e009188b6f4c036c704b086) Fixed #418 -- handle `deis keys:add ~/.ssh/mykey.pub` properly. (@mboersma)
- [5a38afd](https://github.com/deis/deis/commit/5a38afd7c34d39188110f3b22b0878c660eda52f) Added several tests and improved coverage definition. (@mboersma)
- [#420](https://github.com/deis/deis/pull/420) Merge pull request #420 from opdemand/fix-cli-keys-add (@opdemand)
- [#421](https://github.com/deis/deis/pull/421) Merge pull request #421 from opdemand/more-better-tests (@opdemand)
- [95d750f](https://github.com/deis/deis/commit/95d750f502ccc6338af7e954f37b48148379d4f4) Fixed #355 -- retry deleting EC2 security group. (@mboersma)
- [#422](https://github.com/deis/deis/pull/422) Merge pull request #422 from opdemand/fix-delete-ec2-sg (@opdemand)
- [711eba4](https://github.com/deis/deis/commit/711eba4e09183eeac3fdd2adf0e6c90f42f40213) Fixed #402 -- hide web signup, point users to CLI. (@mboersma)
- [#423](https://github.com/deis/deis/pull/423) Merge pull request #423 from opdemand/hide-web-signup (@opdemand)
- [8e5778d](https://github.com/deis/deis/commit/8e5778d5024384db6b2af3b17a9fff08523718a6) Fixed #406 -- return error detail if sudo fails in knife bootstrap. (@mboersma)
- [#425](https://github.com/deis/deis/pull/425) Merge pull request #425 from opdemand/406-knife-bootstrap-err (@opdemand)
- [24f4d0f](https://github.com/deis/deis/commit/24f4d0f8a7310e6f86d6baaafc679f44d15f0a9b) Added tests pointed out by coverage.py. (@mboersma)
- [#426](https://github.com/deis/deis/pull/426) Merge pull request #426 from opdemand/more-tests (@opdemand)
- [99f33ca](https://github.com/deis/deis/commit/99f33cac3c6654d253e70864ac681e9042d9a98c) Removed autofunction docs for non-existent methods. (@mboersma)
- [bf993a2](https://github.com/deis/deis/commit/bf993a2fd90c6878c4219363438cd0e89277e2d3) Fixed #353 -- added permalinks to Sphinx documentation. (@mboersma)
- [#427](https://github.com/deis/deis/pull/427) Merge pull request #427 from opdemand/353-doc-permalinks (@opdemand)
- [a862b91](https://github.com/deis/deis/commit/a862b91050dce4ea7c8f13dd844f303c8e8d56aa) Note that deis only supports `git push` to master, refs #419. (@mboersma)
- [59a4f8b](https://github.com/deis/deis/commit/59a4f8b011091fd5591a3535356812249c5ae860) Note that `nodes:create` requires password-less sudo, refs #362. (@mboersma)
- [#428](https://github.com/deis/deis/pull/428) Merge pull request #428 from opdemand/419-git-push-master-doc (@opdemand)
- [#429](https://github.com/deis/deis/pull/429) Merge pull request #429 from opdemand/362-manual-ssh-sudo (@opdemand)
- [b725bce](https://github.com/deis/deis/commit/b725bce802319b9b94f4ea16ef45c889e2ae97c3) Fixed #430 -- updated EC2 AMIs for 0.3.1 release. (@mboersma)
- [dbae8a2](https://github.com/deis/deis/commit/dbae8a2aa3d56f3178590813ff7c28a41366d8db) Updated release procedure docs. (@mboersma)

### v0.3.0 (2013/12/12 22:06 +00:00)
- [79e4899](https://github.com/deis/deis/commit/79e4899524f1e71b1e3c5916652738c756910761) Switch master to v0.3.0 (@mboersma)
- [c6835f4](https://github.com/deis/deis/commit/c6835f41eedd29dfc8c1b2090219eb31db4477f6) Updated gevent to 1.0, removed Cython dependency. (@mboersma)
- [946c5a7](https://github.com/deis/deis/commit/946c5a765148af11ba5b2de64793327e2d028ed5) add flavors:update to `deis help flavors` (@gabrtv)
- [#341](https://github.com/deis/deis/pull/341) Merge pull request #341 from opdemand/339-gevent-1.0 (@opdemand)
- [25a1f6f](https://github.com/deis/deis/commit/25a1f6ff959f4104b90fe802a168808cb65b2e53) Updated to Django 1.6. (@mboersma)
- [#340](https://github.com/deis/deis/pull/340) Merge pull request #340 from opdemand/django-1.6 (@opdemand)
- [396e3f0](https://github.com/deis/deis/commit/396e3f0e9f84a3e6c74a80685bde00fbf49a1b4f) Fixed query in several build_layer() implementations. (@mboersma)
- [#345](https://github.com/deis/deis/pull/345) Merge pull request #345 from opdemand/344-layer-query-fix (@opdemand)
- [4331cac](https://github.com/deis/deis/commit/4331cac92e3bc94b7f4954208ce76f215fadce0a) Removed the unused allauth.socialaccount app. (@mboersma)
- [2c2e4d9](https://github.com/deis/deis/commit/2c2e4d9cee1e7f9fe2b8b4e84ab4fb624fb7373e) Updated boto, pyrax, django-allauth, and smartypants. (@mboersma)
- [#349](https://github.com/deis/deis/pull/349) Merge pull request #349 from opdemand/remove-socialaccount (@opdemand)
- [#351](https://github.com/deis/deis/pull/351) Merge pull request #351 from opdemand/package-updates (@opdemand)
- [7bd348a](https://github.com/deis/deis/commit/7bd348a7b8c6facbac97c4f5176a9e97bc41f66b) Made vagrant destroy_node() RuntimeError into a warning. (@mboersma)
- [#347](https://github.com/deis/deis/pull/347) Merge pull request #347 from opdemand/346-vagrant-node-dir-err (@opdemand)
- [463dfcf](https://github.com/deis/deis/commit/463dfcfcac9fed865b3cced58411bc951c65b680) Updated to Celery 3.1.6. (@mboersma)
- [#350](https://github.com/deis/deis/pull/350) Merge pull request #350 from opdemand/update-celery (@opdemand)
- [f04cf65](https://github.com/deis/deis/commit/f04cf65cdddd98beabd65bf6903c9a3b7904e312) Fixed missing SECRET_KEY issue for docs generation. (@mboersma)
- [776b425](https://github.com/deis/deis/commit/776b4256130279a77442dbf25549f3e6fdd720af) implement default formation lookup on application creation, with tests (@gabrtv)
- [#357](https://github.com/deis/deis/pull/357) Merge pull request #357 from opdemand/default-formation (@opdemand)
- [7e57b37](https://github.com/deis/deis/commit/7e57b3799a84f1c5f38a19dddbeb8c5031c2bff3) Specify celery worker concurrency on the command-line. (@mboersma)
- [3fb232d](https://github.com/deis/deis/commit/3fb232de00a281639070e54ac00e9d2bcc211d32) downgrade to berkshelf stable, since 3.0.0 beta is problematic (@gabrtv)
- [#358](https://github.com/deis/deis/pull/358) Merge pull request #358 from opdemand/berkshelf-stable (@opdemand)
- [48ac3ee](https://github.com/deis/deis/commit/48ac3ee52d90b35292420868179ee6f9969d54a2) fix celery deadlocks by removing all blocking tasks from other (parent) tasks, and moving the logic to inline model methods (@gabrtv)
- [#359](https://github.com/deis/deis/pull/359) Merge pull request #359 from opdemand/fix-celery-deadlock (@opdemand)
- [b508280](https://github.com/deis/deis/commit/b508280241a706f3254871a2c769f0ca98a7081e) hardcode celeryd concurrency (@gabrtv)
- [5ba327e](https://github.com/deis/deis/commit/5ba327e4fc30d400253e00bc84282983da2d4b20) return the body of chef-client even on failed node convergence (@gabrtv)
- [145a5b6](https://github.com/deis/deis/commit/145a5b679113ae9169548f774c0f374de31bb5ef) Docs - Fix typo (@kumavis)
- [#361](https://github.com/deis/deis/pull/361) Merge pull request #361 from kumavis/patch-1 (@kumavis)
- [#360](https://github.com/deis/deis/pull/360) Merge pull request #360 from opdemand/show-node-converge-error (@opdemand)
- [ac26eb7](https://github.com/deis/deis/commit/ac26eb7c9d43ae6d0a3dd920f3df56612338e34c) update App.run for docker 0.7 and slugrunner syntax (@gabrtv)
- [bf3a57d](https://github.com/deis/deis/commit/bf3a57d284489c78c06833356ac0eab34d3ba11f) switch to docker 0.7.1 and pull progrium/cedarish docker image (@gabrtv)
- [362f95b](https://github.com/deis/deis/commit/362f95b2c9482393889632163d370a1567977b54) Fixed #364 - Updated EC2 AMIs for Deis 0.3.0. (@mboersma)
- [#363](https://github.com/deis/deis/pull/363) Merge pull request #363 from opdemand/docker-0.7 (@opdemand)
- [a27500b](https://github.com/deis/deis/commit/a27500bab2988ef4c1817fabded33fee4a1a7840) initial pass at developer guide and operations guide (@gabrtv)
- [11cae2a](https://github.com/deis/deis/commit/11cae2a12457b279a8e147d9428155d4dedd6a68) Merge branch 'master' into tutorial-docs (@mboersma)
- [575c707](https://github.com/deis/deis/commit/575c707409e08fab8eb17cbcfe9c24870f082b92) update master with latest cookbook for 0.3 (@gabrtv)
- [12dbdbd](https://github.com/deis/deis/commit/12dbdbd3da8a42f941c301ec1be354d6b09cb06d) move chef server to private network along w/ other vagrant components (@gabrtv)
- [a7c6356](https://github.com/deis/deis/commit/a7c6356ceed862b5cac3174e37e6381683aecdcd) fix refs to operations guide (@gabrtv)
- [7e5ff6b](https://github.com/deis/deis/commit/7e5ff6b435e0e2bb2323af2d03cf2cc3740367f3) purge old installation docs (@gabrtv)
- [0e77c57](https://github.com/deis/deis/commit/0e77c576f048d6eee3c5ffe2e9d6a615415bddd4) new local development docs deprecating old devsetup docs (@gabrtv)
- [#374](https://github.com/deis/deis/pull/374) Merge pull request #374 from opdemand/localdev-docs (@opdemand)
- [98bd137](https://github.com/deis/deis/commit/98bd137249088f666ff4b1deed0f17a9d58849eb) restrict converge controller to recipe[deis::gitosis] (@gabrtv)
- [#375](https://github.com/deis/deis/pull/375) Merge pull request #375 from opdemand/upgrade-workflow (@opdemand)
- [8b69f4b](https://github.com/deis/deis/commit/8b69f4b2c1142f08b6cb8c6389048284abba7e14) Removed ref to obsolete app_tasks in formation.destroy. (@mboersma)
- [db6f436](https://github.com/deis/deis/commit/db6f436e2a7f4a8c86edf9cfc65089bd5e9a5c98) set devmode to true for Vagrant controller (@gabrtv)
- [3d48e9f](https://github.com/deis/deis/commit/3d48e9f1c66987adb4646cdbc2bf22b656aa45d8) Fixed app deletion when formation is destroyed. (@mboersma)
- [682d1f8](https://github.com/deis/deis/commit/682d1f80ccb3a0c806b15f00da9470a841840b61) Fix markdown formatting issue in EC2 README (@gabrtv)
- [a707311](https://github.com/deis/deis/commit/a7073118241d783a3db2bc1d76efae0ea190ea9a) remove reuse controller tip (@gabrtv)
- [4790335](https://github.com/deis/deis/commit/4790335faf0f525e9e85c01cd23bc36c1c5feca8) on `deis run` add release environment and auto-remove container (@gabrtv)
- [#378](https://github.com/deis/deis/pull/378) Merge pull request #378 from opdemand/example-apps-tests (@opdemand)
- [#379](https://github.com/deis/deis/pull/379) Merge pull request #379 from opdemand/fix-deis-run (@opdemand)
- [ce13287](https://github.com/deis/deis/commit/ce13287ebc29c107ad2509f357b2d89629bbd592) update Berksfile.lock with latest master (@gabrtv)
- [910146f](https://github.com/deis/deis/commit/910146f3ded7346972cb5e60cb5228a382fdec0d) Check for DO credentials in ENV first then fallback to knife.rb. Add -y to Chef deletion commands. (@tombh)
- [c807cc2](https://github.com/deis/deis/commit/c807cc2d97ea74c2741f9b8c4b6e45dfc34d32a3) Fixed typo in vagrant provision-controller.sh (@mboersma)
- [#348](https://github.com/deis/deis/pull/348) Merge pull request #348 from tombh/digital-ocean-provision-error-tolerance (@tombh)
- [fb87ffd](https://github.com/deis/deis/commit/fb87ffd90282b61db8c70499f16fd703b482c2c0) add admonitions testing section to welcome page (to be reverted after styling) (@gabrtv)
- [86ed2b6](https://github.com/deis/deis/commit/86ed2b676114e026481268b33417c1aa17f88e43) Added icons and styling to Sphinx Admonitions (@bengrunfeld)
- [8f4867a](https://github.com/deis/deis/commit/8f4867a40f4b20325eb49c850f2cabe6563d7cbf) remove admonitions testing from main index (@bengrunfeld)
- [#380](https://github.com/deis/deis/pull/380) Merge pull request #380 from opdemand/sphinx-icons (@opdemand)
- [490e8e0](https://github.com/deis/deis/commit/490e8e0bb6e08858b8ea8d7df5c86d0b7503ceea) Implemented app, formation, superuser sharing permissions. (@mboersma)
- [e363b38](https://github.com/deis/deis/commit/e363b384a8d86bc95a7edfd089167c4bd8308b54) Updated permissions tests. (@mboersma)
- [85954c0](https://github.com/deis/deis/commit/85954c0976f67e0fd59c2926150fd4829522efba) Fixed #368 -- add model custom permissions in migration. (@mboersma)
- [2b8d21a](https://github.com/deis/deis/commit/2b8d21a9f699eb1d5d14218808235c0edc0e654c) Fixed #370 -- added CLI tests for app sharing workflow. (@mboersma)
- [5dcd7a5](https://github.com/deis/deis/commit/5dcd7a573d35102022041ade334b31fffe7bf6bf) Fixed #369 -- added migration to drop djcelery & socialaccount. (@mboersma)
- [#381](https://github.com/deis/deis/pull/381) Merge pull request #381 from opdemand/226-formation-sharing (@opdemand)
- [9529fa0](https://github.com/deis/deis/commit/9529fa086ef8f161251f4ab10ad5a226047248a5) Added basic documentation for `deis sharing` commands. (@mboersma)
- [efdec5c](https://github.com/deis/deis/commit/efdec5c7925c37b780d207e2e9c8862ee32f4b82) fix adominition title and type (@gabrtv)
- [#385](https://github.com/deis/deis/pull/385) Merge pull request #385 from opdemand/sharing-docs (@opdemand)
- [030d0de](https://github.com/deis/deis/commit/030d0dea2a1ae86ea4e7c015ec468f4fd40d0677) working integration suite for example apps (@gabrtv)
- [#386](https://github.com/deis/deis/pull/386) Merge pull request #386 from opdemand/fix-example-tests (@opdemand)
- [ad3028d](https://github.com/deis/deis/commit/ad3028df0316957af1d21711a05b90c2bf8cc20b) swallow the error if the node we're destroying no longer exists (@gabrtv)
- [#387](https://github.com/deis/deis/pull/387) Merge pull request #387 from opdemand/ignore-missing-instance (@opdemand)
- [00f4450](https://github.com/deis/deis/commit/00f4450528760dae2df72603f2ccdd326680004c) Added unit tests for Deis controller web views. (@mboersma)
- [37ce448](https://github.com/deis/deis/commit/37ce448e2854783844c576fd244ad7caeb5dc476) first pass at formation permissions (@gabrtv)
- [7b236d8](https://github.com/deis/deis/commit/7b236d8513dc7a33dcc4fb56201e858e0ce3dca2) Removed the unused /docs/ web view. (@mboersma)
- [fe4af27](https://github.com/deis/deis/commit/fe4af27488718bbce685f0b46ae364a84e590532) raise ResponseError on perms operations (@gabrtv)
- [#388](https://github.com/deis/deis/pull/388) Merge pull request #388 from opdemand/web-view-tests (@opdemand)
- [#389](https://github.com/deis/deis/pull/389) Merge pull request #389 from opdemand/formations-perms (@opdemand)
- [8fa1014](https://github.com/deis/deis/commit/8fa1014787b23a6452831bd47ce9cd2caac6cd62) Added is_staff flag to initial superuser, refs #389. (@mboersma)

### v0.2.1 (2013/11/26 18:41 +00:00)
- [5f8d83b](https://github.com/deis/deis/commit/5f8d83be50767bdf350856f09f2baf01b05a8247) Switch master to v0.2.1 (@mboersma)
- [ce2afc6](https://github.com/deis/deis/commit/ce2afc6a11edbb6c237550e6afd55d85178a4868) Started on docs (@mboersma)
- [f69b0da](https://github.com/deis/deis/commit/f69b0daf01c03910c1eb06dec31fe3f1564aaf21) Fixed Flake8 errors (@tombh)
- [03ed368](https://github.com/deis/deis/commit/03ed3681a70a31e8fd1b426606d59910b7e6d8af) Updated Berksfile format for 3.0.0 beta 3 and Chef to 11.6.2. (@mboersma)
- [8fff79e](https://github.com/deis/deis/commit/8fff79e30474df5e3d920ade041b28a5ca87788f) Silence the McCabe checker on two methods, refs #276. (@mboersma)
- [e86d7c3](https://github.com/deis/deis/commit/e86d7c35574ec0fcd463ad95f63788663b093634) Changed sync_folder path in Vagrantfile.local.example to mount parent path (@tombh)
- [34c70c8](https://github.com/deis/deis/commit/34c70c81368334be9f7fbfb8df7894192ac78f5b) Full Vagrant provider. Squashed commits (@tombh)
- [17efd9d](https://github.com/deis/deis/commit/17efd9dc75168109f7946c0740d6cad6eff59ae6) Added missing warning to vagrant provisioning script. (@mboersma)
- [4786ec8](https://github.com/deis/deis/commit/4786ec8c439dd53256413a58644b6d0aef939a55) Default login/register URL to http:// schema if it's not specified. (@mboersma)
- [7ecb90f](https://github.com/deis/deis/commit/7ecb90f77d4d380bae6a4a8b35b13f98fc955a8e) Moved "Terms" section back to top level nav in docs. (@mboersma)
- [#300](https://github.com/deis/deis/pull/300) Merge pull request #300 from opdemand/299-no-schema-supplied (@opdemand)
- [#301](https://github.com/deis/deis/pull/301) Merge pull request #301 from opdemand/296-terms-sphinx-nav (@opdemand)
- [13d0adc](https://github.com/deis/deis/commit/13d0adcc622331e1c7c63668fc0771db59fe4553) Don't prompt users to create an app if runtime=0 (@mboersma)
- [#302](https://github.com/deis/deis/pull/302) Merge pull request #302 from opdemand/250-nodes-scale-msg (@opdemand)
- [28e395a](https://github.com/deis/deis/commit/28e395a6ea5001add742e92dd08c261765e56132) Experimenting with using hostname rather than IP to SSH into host. (@tombh)
- [96ae777](https://github.com/deis/deis/commit/96ae777e9bb8de4012ecf0e5c1cd2d7b80d5f435) 1) Use private network and static IPs for all vagrant VMs. 2) Use (@tombh)
- [bf234f7](https://github.com/deis/deis/commit/bf234f72d07962e4a9ad9db0185b05613d48d0f2) Override rsyslog behaviour by using Vagrantfile to create a higher (@tombh)
- [6083875](https://github.com/deis/deis/commit/6083875b8ef0a4cd5b968627ac100d94a4595958) added some error checking for digital ocean
- [4c74f1e](https://github.com/deis/deis/commit/4c74f1eeb621a1119f4e5e6701b2e0df08ef145d) use region_id, not location_id
- [6ef32c7](https://github.com/deis/deis/commit/6ef32c71986b3e6ed825b457a0f430955d7a7850) Refactored provider creds discovery to be more DRY, added alternate AWS vars. (@mboersma)
- [#308](https://github.com/deis/deis/pull/308) Merge pull request #308 from bacongobbler/306-snapshot-error-checking (@bacongobbler)
- [#309](https://github.com/deis/deis/pull/309) Merge pull request #309 from opdemand/276-providers-discover-refactor (@opdemand)
- [f0fd45f](https://github.com/deis/deis/commit/f0fd45f398641266e40c1081e7ac0920684b42fc) If CM's purge_node() fails then raise an error. (@tombh)
- [0c0f7ca](https://github.com/deis/deis/commit/0c0f7caa73e5d7fb4be4000ad8a3d6cf32f23172) If CM's purge_node() fails then raise an error. (@tombh)
- [#310](https://github.com/deis/deis/pull/310) Merge pull request #310 from tombh/purge-node-feedback (@tombh)
- [be45748](https://github.com/deis/deis/commit/be457480450a26ec152be7a494b903904c8e315b) Instructions to add Controller to Chef's admin group and note to install avahi-daemon/Bonjour (@tombh)
- [691461d](https://github.com/deis/deis/commit/691461dc85ecca3e61334e06e239fda93c890e88) Merge branch '232-vagrant-provider-full' of https://github.com/tombh/deis into tombh-232-vagrant-provider-full (@mboersma)
- [66493b1](https://github.com/deis/deis/commit/66493b1058b53939cf415d23cdba2bf90b5675f9) Updated docs and only test for avahi-daemon on Linux. (@mboersma)
- [3e3c14a](https://github.com/deis/deis/commit/3e3c14a9b9e7d9cfc5c456655ee3a6c1251dafe4) Merge branch 'patch-2' of https://github.com/scottstamp/deis into scottstamp-patch-2 (@mboersma)
- [5d5bf1b](https://github.com/deis/deis/commit/5d5bf1b6080073c48ab8ee3351a559c6a252421c) Updated Berksfile for latest deis-cookbook SHA. (@mboersma)
- [f18f32c](https://github.com/deis/deis/commit/f18f32cf4ad32cb0ff10b85f05004264b2c16148) Made vagrant provider's config not a fatal IOError. (@mboersma)
- [#314](https://github.com/deis/deis/pull/314) Merge pull request #314 from opdemand/313-vagrant-ioerror (@opdemand)
- [21d0634](https://github.com/deis/deis/commit/21d0634a7ea84d25a9691a82d39c227bfe2e2f45) Ignore Chef 404s when destroying a node. (@mboersma)
- [#315](https://github.com/deis/deis/pull/315) Merge pull request #315 from opdemand/312-node-purge-404 (@opdemand)
- [ef4b021](https://github.com/deis/deis/commit/ef4b021e3807ff61808cd57fdfbb0b039118c478) Created a separate pip requirements file for doc generation. (@mboersma)
- [#316](https://github.com/deis/deis/pull/316) Merge pull request #316 from opdemand/311-orphaned-sphinx-docs (@opdemand)
- [b59f28b](https://github.com/deis/deis/commit/b59f28b0c290b26dc2d14487ad7e5f37c7779dc8) Added more helpful error when `deis run` comes before `git push`. (@mboersma)
- [#323](https://github.com/deis/deis/pull/323) Merge pull request #323 from opdemand/304-run-before-push (@opdemand)
- [19b6e5b](https://github.com/deis/deis/commit/19b6e5be39d430195db02f48a96a7793860bebc8) Added ssh-key generation per test user. (@mboersma)
- [#324](https://github.com/deis/deis/pull/324) Merge pull request #324 from opdemand/214=cli-acceptance-tests (@opdemand)
- [503f282](https://github.com/deis/deis/commit/503f28205d269f6d28d09fcbc5660dfc1293e5b3) Create fake home dir per test user, refs #214. (@mboersma)
- [#330](https://github.com/deis/deis/pull/330) Merge pull request #330 from opdemand/214-cli-acceptance-tests (@opdemand)
- [97781ee](https://github.com/deis/deis/commit/97781ee8310486e4f155359870f0963e5ca71c8b) Added test_examples to hit each example-* project, WIP refs #214. (@mboersma)
- [8ce7fa0](https://github.com/deis/deis/commit/8ce7fa005099085f078c8d9c83f7c4665172cdac) Implemented missing chef.purge_user() and connected it to a signal. (@mboersma)
- [#332](https://github.com/deis/deis/pull/332) Merge pull request #332 from opdemand/331-cm-purge-user (@opdemand)
- [d4c1938](https://github.com/deis/deis/commit/d4c1938e660c101f4fab3f1f024c3547c9bf67b3) Added documentation for static / bare metal installation. (@mboersma)
- [#333](https://github.com/deis/deis/pull/333) Merge pull request #333 from opdemand/214-cli-tests (@opdemand)
- [#335](https://github.com/deis/deis/pull/335) Merge pull request #335 from opdemand/295-static-installation (@opdemand)
- [f34ea7d](https://github.com/deis/deis/commit/f34ea7d70fb39588013355b9878c6a81d1529e08) Fixed PHP detect regex in test_examples, refs #214. (@mboersma)
- [c75566d](https://github.com/deis/deis/commit/c75566d14386063c39184e178068f124f75492ad) Increase vagrant controller RAM to 2G, closes #336. (@mboersma)
- [6756763](https://github.com/deis/deis/commit/67567630ba4f6109757b3ba5df6f7191b6a3d0a8) fix jquery in airplane-mode (@gabrtv)
- [7c4baee](https://github.com/deis/deis/commit/7c4baeea9120b20b041d25105e4fafcbd38c9337) lay down initial structure dev/ops tutorials (@gabrtv)
- [f367c7e](https://github.com/deis/deis/commit/f367c7e26838672892a8d9e1e63340d92c304fc2) Fixed #327 -- updated EC2 AMIs for 0.2.1 release. (@mboersma)

### v0.2.0 (2013/11/05 18:56 +00:00)
- [d783c24](https://github.com/deis/deis/commit/d783c244db85a7b152c4578f3d208435f9399001) Switch master to v0.1.2 (@mboersma)
- [29929dd](https://github.com/deis/deis/commit/29929dd507b5152a59fba0b19a6ee2f99bacb4f3) Fixed reference to "only EC2" in client README.rst. (@mboersma)
- [5f83a5e](https://github.com/deis/deis/commit/5f83a5e87cd0ea6702b10a2a9826a3cf99ec6a72) Updated Django to 1.5.5 security/bugfix release. (@mboersma)
- [4190cef](https://github.com/deis/deis/commit/4190cef097846f470c3387544f44a0d7458ec8a4) Updated opdemand/buildstep build procedure. (@mboersma)
- [0c6f8cd](https://github.com/deis/deis/commit/0c6f8cd2c08d2db047b56283dec2378e7d8f2443) Fixed #248 -- exposed `deis layers:update` in CLI. (@mboersma)
- [7f17ffb](https://github.com/deis/deis/commit/7f17ffbc6ebdd6a89440f343e36f3cb3ddce1339) Fixed #242 -- update "no proxy" error message. (@mboersma)
- [#253](https://github.com/deis/deis/pull/253) Merge pull request #253 from opdemand/248-layers-update (@opdemand)
- [c9127c1](https://github.com/deis/deis/commit/c9127c1b80a666c1e9bf86d3a6d31cfb0a384c8b) Fixed TypeError in formations:create. (@mboersma)
- [5bcae3f](https://github.com/deis/deis/commit/5bcae3f20b2a84f7a4cc7e2328f650790ec47f88) Fixed #232 -- add vagrant support for Deis development. (@mboersma)
- [3b45465](https://github.com/deis/deis/commit/3b45465c5f7dfbfc4089d5dae939a1e9276b5624) validate that App.id only contains [a-z0-9-] and is a valid domain name, with tests (@gabrtv)
- [94d714d](https://github.com/deis/deis/commit/94d714df444efbf2224e5b52d8ed3007c24f7462) fix variable naming to be clear about app vs. formation (@gabrtv)
- [1ca8d84](https://github.com/deis/deis/commit/1ca8d840f5b655a269370aa9d7a8fe8e7c8da577) add digital ocean provider
- [878e769](https://github.com/deis/deis/commit/878e769e7dab9076e3271d6605af93729af683b7) Added a script to help create the vagrant static formation. (@mboersma)
- [#254](https://github.com/deis/deis/pull/254) Merge pull request #254 from opdemand/fix-id-override (@opdemand)
- [#257](https://github.com/deis/deis/pull/257) Merge pull request #257 from opdemand/232-vagrant-dev (@opdemand)
- [3ba4960](https://github.com/deis/deis/commit/3ba4960fc323a30aef862ce19bade60d672743c4) changed controller size to 2GB
- [e054ffb](https://github.com/deis/deis/commit/e054ffb9674392ede1b8a22b9ff3b59ed855639a) Merge branch '73-digitalocean-provider' of https://github.com/bacongobbler/deis into bacongobbler-73-digitalocean-provider (@mboersma)
- [765a9cc](https://github.com/deis/deis/commit/765a9ccf98b8515c7627d46a87845df94410a142) Fixed PEP8 errors missed on previous merge. (@mboersma)
- [124ff26](https://github.com/deis/deis/commit/124ff26914cec56f56da485e0ac5bf70074389f1) Include provider.digitalocean in API docs, refs #73. (@mboersma)
- [4e4303a](https://github.com/deis/deis/commit/4e4303aba10031ea80059136991c27cf0431974f) Fixed #258 -- removed errant print statement. (@mboersma)
- [ce064c5](https://github.com/deis/deis/commit/ce064c56f0a59d9836291c8f7b1149a9eed568b4) remove ffi dependency
- [#261](https://github.com/deis/deis/pull/261) Merge pull request #261 from bacongobbler/patch-1 (@bacongobbler)
- [572a1c9](https://github.com/deis/deis/commit/572a1c9115f832beae9b8fcb994a4b6f6f9f9bcc) downgrade eventmachine version to v1.0.0
- [#263](https://github.com/deis/deis/pull/263) Merge pull request #263 from bacongobbler/patch-2 (@bacongobbler)
- [7740d7f](https://github.com/deis/deis/commit/7740d7f9247f45f89fb0bd7a20b9bd0cd6ea09e7) remove EC2 only from README (@gabrtv)
- [4cf312f](https://github.com/deis/deis/commit/4cf312f3505b23422c5a6d2847f6c38babf0e6d0) Proposed fix for gh#264 (@scottstamp)
- [b666811](https://github.com/deis/deis/commit/b666811335cb5e3baa657b7017a4440ef05d85e0) Update chef.py (@scottstamp)
- [#265](https://github.com/deis/deis/pull/265) Merge pull request #265 from scottstamp/patch-1 (@scottstamp)
- [39e211c](https://github.com/deis/deis/commit/39e211c395ff91e76799c0a5b9ac300c88e6f9a3) Added a reminder about Chef admin perms after provisioning. (@mboersma)
- [c132663](https://github.com/deis/deis/commit/c1326630a6b4c7405d96c8b65e3b7ca6957d6acd) Use correct method in socket module, gethostname() (@gabrtv)
- [b800549](https://github.com/deis/deis/commit/b800549dd8ea0479f5c04bf56abc94396ffb6bbb) Fixed #267 -- updated controller web UI. (@mboersma)
- [#272](https://github.com/deis/deis/pull/272) Merge pull request #272 from opdemand/267-web-ui (@opdemand)
- [b8d2cbf](https://github.com/deis/deis/commit/b8d2cbf7b9d549e5ab214196bae8b51d2314f735) Updated Chef version to 11.6.2. (@mboersma)
- [358d36b](https://github.com/deis/deis/commit/358d36b15f0d758d29b2c1ff5ed9b989a700cf9d) Updated to new Berksfile format to silence deprecation warnings. (@mboersma)
- [c72e970](https://github.com/deis/deis/commit/c72e9708e23b912bb33590cfdf651248d59f6894) Created provider-specific installation docs, refs #227. (@mboersma)
- [#283](https://github.com/deis/deis/pull/283) Merge pull request #283 from opdemand/227-provider-docs (@opdemand)
- [#280](https://github.com/deis/deis/pull/280) Merge pull request #280 from opdemand/271-chef-version (@opdemand)
- [059ba16](https://github.com/deis/deis/commit/059ba16b5d5294e1e2c23018c4a030dbefd28c3e) Reverted Berksfile format. (@mboersma)
- [5ee80eb](https://github.com/deis/deis/commit/5ee80ebf6e4f310a7323a012cfac03e1c280ac1a) Updated sadly hard-coded JS nav code for new docs, refs #283. (@mboersma)
- [1eb1dc4](https://github.com/deis/deis/commit/1eb1dc47f3bc4eedacf31ab1187173219da70748) Clarify that Dockerfiles are not yet supported directly. (@mboersma)
- [1838ae3](https://github.com/deis/deis/commit/1838ae3e1014bcf00b36ef0f4895e25d38223d1f) Updated several python packages. (@mboersma)
- [29ba038](https://github.com/deis/deis/commit/29ba038207219c33efded90bc700ab80f77a995f) Specified package version 0.6.4 for Docker, not virtual package. (@mboersma)
- [b4b1272](https://github.com/deis/deis/commit/b4b127259b9b56f877a3537a7b16e257ec5d6aaa) Updated EC2 AMIs in all regions for security updates. (@mboersma)
- [#289](https://github.com/deis/deis/pull/289) Merge pull request #289 from opdemand/236-refresh-amis (@opdemand)
- [cb97507](https://github.com/deis/deis/commit/cb97507703a8ee320132a411e757ebfca59f21f3) Updated prepare-image-* scripts for minor optimizations. (@mboersma)
- [ef9290b](https://github.com/deis/deis/commit/ef9290b81590db602022dad407a5c9dda671c4f6) Updated project version strings to 0.2.0. (@mboersma)
- [d25c2c2](https://github.com/deis/deis/commit/d25c2c2c233f7001a42de3680c8fbe8128f4eff9) Updated Ruby gems dependencies. (@mboersma)

### v0.1.1 (2013/10/22 16:42 +00:00)
- [4836b99](https://github.com/deis/deis/commit/4836b998fb6bb627f00bd58093d33e7414b20fe2) switch master to v0.1.1 (@gabrtv)
- [3c64f7f](https://github.com/deis/deis/commit/3c64f7fb34f7fcbe0d94425ce42d372b9375d7e2) update release instructions post git-flow (@gabrtv)
- [12d6f07](https://github.com/deis/deis/commit/12d6f07f98ebbcafb4c7a9999d76fc8a5bf1acbd) Fixed a RST format bug in client README.rst. (@mboersma)
- [2228d6c](https://github.com/deis/deis/commit/2228d6ca16f02e29f9520128149d2e9be9c3b751) Fixed #215 -- allow flavors:update from CLI. (@mboersma)
- [#217](https://github.com/deis/deis/pull/217) Merge pull request #217 from opdemand/215-flavors-update (@opdemand)
- [f4f3ae3](https://github.com/deis/deis/commit/f4f3ae38c5cd6d3e91c07aa02d3ba4630c097b9c) Added bing webmaster tools verification meta (@bengrunfeld)
- [#218](https://github.com/deis/deis/pull/218) Merge pull request #218 from opdemand/add-bing-meta (@opdemand)
- [f3154e6](https://github.com/deis/deis/commit/f3154e6196f3dd83200dd52ff032450af6cebd7c) Updated client version to 0.1.1 per release process. (@mboersma)
- [767f090](https://github.com/deis/deis/commit/767f090fe8ddf19fe5a93ef3ca9aa0294fa6eec1) Remove refs to old 'azure' attempt. See also #219. (@mboersma)
- [1746bc4](https://github.com/deis/deis/commit/1746bc4c385d468ee8c13ee8fab2ce56ef13fab9) Updated paramiko and Sphinx. (@mboersma)
- [f320a77](https://github.com/deis/deis/commit/f320a7740000a2f198ced6f4b6884f4825ec55e0) Updated coverage to 3.7. (@mboersma)
- [b848ae3](https://github.com/deis/deis/commit/b848ae37f02280c8261e626ec76eb68d157e2b0e) Fixed #121 -- added support for Rackspace open cloud (@mboersma)
- [5c00f05](https://github.com/deis/deis/commit/5c00f057bbe46a086d7e6671a780db6ef45578c4) Fixed a few typos. (@mboersma)
- [9605456](https://github.com/deis/deis/commit/9605456da0522a6eb56fb766d4b3c9953bacd47e) update docs sitemap.xml post 0.1.0 refactoring (@gabrtv)
- [#224](https://github.com/deis/deis/pull/224) Merge pull request #224 from opdemand/121-rackspace-provider (@opdemand)
- [#225](https://github.com/deis/deis/pull/225) Merge pull request #225 from opdemand/update-sitemap (@opdemand)
- [35c371f](https://github.com/deis/deis/commit/35c371ff9df1706f83541665773265a4c2502a72) Fixed #222 -- CLI scaling error msg when no Rackspace image exists (@mboersma)
- [#228](https://github.com/deis/deis/pull/228) Merge pull request #228 from opdemand/222-better-500-error (@opdemand)
- [217f7dc](https://github.com/deis/deis/commit/217f7dcc22b71394d0053c228bba33c850ee1e16) only run commands against nodes in runtime layers (@gabrtv)
- [#229](https://github.com/deis/deis/pull/229) Merge pull request #229 from opdemand/fix-run-selection (@opdemand)
- [551e742](https://github.com/deis/deis/commit/551e742fa727a4523e44caf9b08fc6fb18acee80) Fixed possible exception case from previous commit. (@mboersma)
- [77a571f](https://github.com/deis/deis/commit/77a571fdf46274531428c7568afd59a7ce830412) Fixed #233 -- made buildstep fork cleaner WRT progrium/buildstep (@mboersma)
- [05ff02a](https://github.com/deis/deis/commit/05ff02a12461bd61a71e7b3ae6a0b0ec9aa636ba) Fixed #220 -- add standard footer to Sphinx template for readthedocs.org (@mboersma)
- [6b1a5b3](https://github.com/deis/deis/commit/6b1a5b3175c476be2c4ab546ebdf2a63d4abe3b5) Updated boto, pycrypto, yamlfield. (@mboersma)
- [b58105a](https://github.com/deis/deis/commit/b58105af2207b3fc08b9a36b61f7db6b0dcbef90) Updated gevent to 1.0rc3, fixes vagrant DNS issues, refs #232. (@mboersma)
- [278244e](https://github.com/deis/deis/commit/278244eb6c55e5c74b9a75a7e72e388bb50dcaab) Fixes #240. Updated docs layout.html file and main.css to fix the style breakage that adding Read The Docs footer caused. (@bengrunfeld)
- [85444c0](https://github.com/deis/deis/commit/85444c069b2a5ecafed9fd012d4d4591c5b1e8fb) Updates main.css to fix styling break caused by ReadTheDocs footer template tag. (@bengrunfeld)
- [d6eeeec](https://github.com/deis/deis/commit/d6eeeec8b53aecc8741a53d1cd95b9c5e3838dc7) refs #240. Hopefully fixes style breakages. (@bengrunfeld)
- [5fc6d3d](https://github.com/deis/deis/commit/5fc6d3debb78e25b151b2dd18bc8312d018ec5a0) refs #240. Hopefully fixes style breakages. (@bengrunfeld)
- [c0d6576](https://github.com/deis/deis/commit/c0d6576f9bb920620abd1a71ef4776a422c0f484) Fixed #234 -- refreshed EC2 AMI images. (@mboersma)
- [#241](https://github.com/deis/deis/pull/241) Merge pull request #241 from opdemand/234-refresh-amis (@opdemand)
- [ec4899c](https://github.com/deis/deis/commit/ec4899c9830db59d9e91193348d59abdd1b952c2) Minor update to release procedure doc. (@mboersma)

### v0.1.0 (2013/10/01 22:14 +00:00)
- [e027393](https://github.com/deis/deis/commit/e027393920181373c216864aa1aca33c79db7927) switch master to v0.0.8 (@gabrtv)
- [a24f0dd](https://github.com/deis/deis/commit/a24f0dd59d772d2d5255994f278a996a0e79f361) Updated release procedure notes. (@mboersma)
- [bc27e99](https://github.com/deis/deis/commit/bc27e9969cf63d765f7c670201a9fb4cd8f0eace) refactor wip (@gabrtv)
- [5280796](https://github.com/deis/deis/commit/52807965b84f431622645fa4e0616074cdcb8330) Updated to Django 1.5.3 security release. (@mboersma)
- [032f908](https://github.com/deis/deis/commit/032f908312cb9c8f853954def955043eaeafa14d) updated to Django v1.5.4 security release
- [#173](https://github.com/deis/deis/pull/173) Merge pull request #173 from bacongobbler/172-django-174-security-release (@bacongobbler)
- [df50eee](https://github.com/deis/deis/commit/df50eee1aeb3b478bcbe1cb5b56c1da8a614aec1) add test coverage for cm synchronization using mock CM module (@gabrtv)
- [0edb429](https://github.com/deis/deis/commit/0edb429825b2c10fdf099e879373e31260f1314f) create the 3 necessary data bags (@gabrtv)
- [3f33765](https://github.com/deis/deis/commit/3f33765701eba1bd04a70977f01f87d71c861ffe) fix bug in key deletion api call (@gabrtv)
- [f306aad](https://github.com/deis/deis/commit/f306aadc39d7258b634730da1214a39b5b9ff3f8) move run_node to CM module, add container list/info endpoints, beef test coverage (@gabrtv)
- [74dfe2e](https://github.com/deis/deis/commit/74dfe2e6264dffb3ceba2094390f6d41889c8f58) check for initial web container in build.push (@gabrtv)
- [d045377](https://github.com/deis/deis/commit/d045377c14b62c01d7b3b968bead3d513c24ab34) add Build.push test coverage (@gabrtv)
- [591ba16](https://github.com/deis/deis/commit/591ba16ab5ee467d5d1d0f10c15bfd520e196df4) fix layer tests (@gabrtv)
- [ff28545](https://github.com/deis/deis/commit/ff2854532e55ed0b050a126097e0231e3b5c3d41) converge on app scale operation (@gabrtv)
- [ea24ebb](https://github.com/deis/deis/commit/ea24ebb729f9aa625de888479e3432e6a048d728) Updated to Django 1.5.3 security release. (@mboersma)
- [8abeea0](https://github.com/deis/deis/commit/8abeea03918110921c75719d6ec43b0861ec0ecd) updated to Django v1.5.4 security release
- [60bdbdd](https://github.com/deis/deis/commit/60bdbdde6642519d73699cc8d0530f95e70e84f5) improve code coverage on api views/models (@gabrtv)
- [f62a0cb](https://github.com/deis/deis/commit/f62a0cb31d9b09b556962dccd2f504adeb68de49) add apps:calculate functionality (@gabrtv)
- [7f8debc](https://github.com/deis/deis/commit/7f8debcd9513cf92a8426ed17ee158910f81d19d) add containers to application databag (@gabrtv)
- [7ae9c90](https://github.com/deis/deis/commit/7ae9c9050096b26f3f0d40df96c00a22a5e6aaa0) add nodes:converge, nodes:ssh, fix `deis open` (@gabrtv)
- [f88e782](https://github.com/deis/deis/commit/f88e78286a711bf3f2146674564e37bf39ccf720) move logs to apps:logs and create shortcut (@gabrtv)
- [e6b57a0](https://github.com/deis/deis/commit/e6b57a0cd6a919f5693a6fa7783a7e5c1e4c1149) fix bug with git push working after formation destroy (@gabrtv)
- [77e266a](https://github.com/deis/deis/commit/77e266a37598561ada8b322f9cb555322f6e0dda) switch release endpoints to app (@gabrtv)
- [df72bc7](https://github.com/deis/deis/commit/df72bc7dcc8ff6772442d7d0f2edd907d09155ce) change config:list, set and unset to use app endpoints (@gabrtv)
- [eb7067f](https://github.com/deis/deis/commit/eb7067faa7486b566a97e963a9a73b08d6ef9c22) move open to apps:open (@gabrtv)
- [a05b5c8](https://github.com/deis/deis/commit/a05b5c89c70d0d23278bc192e623e79a145f5e68) change app databag format from proxies to domains (@gabrtv)
- [c64c2e2](https://github.com/deis/deis/commit/c64c2e2bcfb271b19e17901bdf303b682ed16667) add test coverage around mutiple apps for formations that do/dont support it (@gabrtv)
- [28af652](https://github.com/deis/deis/commit/28af65204e2115eae52c4e3913815b1c81ab16ea) move image into build for better 12 factor compatibility (@gabrtv)
- [2e0eb69](https://github.com/deis/deis/commit/2e0eb6920fca13ffc6fdd68ba64879c6770e2ce5) add formations:update to set domain used in multi-app (@gabrtv)
- [318d9b4](https://github.com/deis/deis/commit/318d9b4468c5a3cd32f0895571434a74e6cd72cf) add container.port and make it unique across a formation for multi-app (@gabrtv)
- [f6f13b4](https://github.com/deis/deis/commit/f6f13b4ac3a8913d4a30f0950550ed04d22eec12) Fixed #174 -- updated Sphinx docs to match refactored modules. (@mboersma)
- [cd28ce4](https://github.com/deis/deis/commit/cd28ce404d140e7055db3abeb2fe19509e7f3a19) remove ops commands from client help (@gabrtv)
- [f16dd02](https://github.com/deis/deis/commit/f16dd0215ce0887bf3d67922a1b73aafb8be3787) Fixed #175 -- update REST API documentation for application-refactor. (@mboersma)
- [3ad53ce](https://github.com/deis/deis/commit/3ad53ce38496c37eb4ccd3fa0cb636cafac2ef9b) fix missing command bug during dispatching (@gabrtv)
- [131bf8f](https://github.com/deis/deis/commit/131bf8fdce1d15bed491336021cbe3f0fae080b1) fix bug with apps:run task dispatch (@gabrtv)
- [e54f9e2](https://github.com/deis/deis/commit/e54f9e2849f0214cb1a24abcd0aed0c1a46b527d) fix formation databag after app destroy, with test (@gabrtv)
- [#179](https://github.com/deis/deis/pull/179) Merge pull request #179 from opdemand/application-refactor (@opdemand)
- [d05ca70](https://github.com/deis/deis/commit/d05ca706964da5b5d20736b6ef5ee8a920107be6) deprecate provider.controller module, now handled by cm package (@gabrtv)
- [3e12255](https://github.com/deis/deis/commit/3e122552ef840151fe48bafa4c700e5807a636de) Removed Sphinx docs referencing obsolete provider.controller module. (@mboersma)
- [#180](https://github.com/deis/deis/pull/180) Merge pull request #180 from opdemand/deprecate-old-providers (@opdemand)
- [f856c67](https://github.com/deis/deis/commit/f856c671adc6485895bd454109f0b62c761a49f8) Fixed #167 -- allow underscore in slug-type field regexes. (@mboersma)
- [#183](https://github.com/deis/deis/pull/183) Merge pull request #183 from opdemand/167-underscores-in-name (@opdemand)
- [a646fbc](https://github.com/deis/deis/commit/a646fbc2ff17ecc05181184ac00bff2f8e443610) add ssh public key to Node.flat() data structure (@gabrtv)
- [#184](https://github.com/deis/deis/pull/184) Merge pull request #184 from opdemand/add-node-ssh-pubkey (@opdemand)
- [99da003](https://github.com/deis/deis/commit/99da003dce661d35b6c31d209438228281d0a7f8) Fixed #176 -- added docstrings for provider modules. (@mboersma)
- [#185](https://github.com/deis/deis/pull/185) Merge pull request #185 from opdemand/176-provider-docstrings (@opdemand)
- [2b4d7f8](https://github.com/deis/deis/commit/2b4d7f8df9ca74b7de32d3693bad18516aed8ddb) Fixed #177 -- added docstrings for cm modules. (@mboersma)
- [#186](https://github.com/deis/deis/pull/186) Merge pull request #186 from opdemand/177-cm-docstrings (@opdemand)
- [4f04098](https://github.com/deis/deis/commit/4f04098005777af0fd553f28561cae4634e05f5d) cli doc updates and cleanup for app refactoring (@gabrtv)
- [a254564](https://github.com/deis/deis/commit/a254564ba79664178ff4eefee9848f3dcff19a9a) update client reference docs (@gabrtv)
- [ca46506](https://github.com/deis/deis/commit/ca46506e04d27cb3229a8975b59564e3ab781f7b) update terms post application object model changes (@gabrtv)
- [5f09ad8](https://github.com/deis/deis/commit/5f09ad8c3c68889f0205eed6ed5cbc3eef64eaca) update concepts post application object model changes (@gabrtv)
- [75487ed](https://github.com/deis/deis/commit/75487ed48c61df6cb9152e73fed7a97ebb255a85) add apps:info with info shortcut (@gabrtv)
- [8940321](https://github.com/deis/deis/commit/89403219dd25bfc08b3e1171a8b4bb826e487dc5) remove unnecesary JSONField subclasses, remote undoc'd members from api.models (@gabrtv)
- [ed22d9a](https://github.com/deis/deis/commit/ed22d9ab73c9d24d9066af82b87c6247db272ca7) standardize on spaces instead of tabs to fix code block formatting (@gabrtv)
- [3d9142f](https://github.com/deis/deis/commit/3d9142f7ade378386037af7be8f07cc19c4257c1) fix noindex indentation (@gabrtv)
- [425d498](https://github.com/deis/deis/commit/425d498798b98143b29f7d8ce4a558c5e9bd4026) add escaping of astericks for sphinx (@gabrtv)
- [66d6427](https://github.com/deis/deis/commit/66d6427e12d5687a4e80611ac9b3aec2a6ee6257) update installation docs (@gabrtv)
- [c7c6e9e](https://github.com/deis/deis/commit/c7c6e9ec60b93282c500d7973aa46d71d8f639e3) remove todos from toctree (@gabrtv)
- [65245e2](https://github.com/deis/deis/commit/65245e2f19b9f29cf5d9c9cb721b1436bb8d1f99) update readme with new app-oriented workflow (@gabrtv)
- [369ac79](https://github.com/deis/deis/commit/369ac795c38e86f9bbe9267f4f6eff50280eea56) fix typo in deploy language (@gabrtv)
- [#189](https://github.com/deis/deis/pull/189) Merge pull request #189 from opdemand/doc-updates (@opdemand)
- [fda0f27](https://github.com/deis/deis/commit/fda0f2706f0a237dfe71251b6730d250ba8baf6b) return 404 if app cannot be found during queryset lookup, fixes #182 (@gabrtv)
- [2f269fa](https://github.com/deis/deis/commit/2f269fad14330e9b24a87071c6aa6569c12b95eb) add simple password confirmation on registration fixes #188 (@gabrtv)
- [#192](https://github.com/deis/deis/pull/192) Merge pull request #192 from opdemand/confirm-pw-on-register (@opdemand)
- [#191](https://github.com/deis/deis/pull/191) Merge pull request #191 from opdemand/app-endpoint-500s (@opdemand)
- [fdcd7d9](https://github.com/deis/deis/commit/fdcd7d9710facf7a10291ee1504239ab45021877) more doc updates (@gabrtv)
- [7e3230b](https://github.com/deis/deis/commit/7e3230b909af21a0d2c2e15ac429fa4bf8a22aa4) cleanup line breaks, update command output (@gabrtv)
- [678a6b8](https://github.com/deis/deis/commit/678a6b897156019afce22aa469a0e46d26bf75d4) reset initial south migration, fixes opdemand/deis-cookbook#14 (@gabrtv)
- [#194](https://github.com/deis/deis/pull/194) Merge pull request #194 from opdemand/more-doc-updates (@opdemand)
- [#193](https://github.com/deis/deis/pull/193) Merge pull request #193 from opdemand/reset-south-migrations (@opdemand)
- [8e4f60f](https://github.com/deis/deis/commit/8e4f60f8a04a52a7642607d7234e35a92bc8bbbb) Updated and tested boto, requests, djangorestframework, and paramiko. (@mboersma)
- [#195](https://github.com/deis/deis/pull/195) Merge pull request #195 from opdemand/pypi_updates (@opdemand)
- [0388366](https://github.com/deis/deis/commit/0388366b3eeb0f3397cbd4ea6462e6bc74fe9a08) fix ssh_username not being read in 'deis ssh'
- [#196](https://github.com/deis/deis/pull/196) Merge pull request #196 from bacongobbler/deis-ssh-username-fix (@bacongobbler)
- [4df6164](https://github.com/deis/deis/commit/4df6164e3084ca5dc7219e5f2812d1124912b18d) Added a django admin class for api.models.App. (@mboersma)
- [#197](https://github.com/deis/deis/pull/197) Merge pull request #197 from opdemand/app-in-django-admin (@opdemand)
- [034a4ac](https://github.com/deis/deis/commit/034a4ac500ad2164a713512555625112b0d3e344) Fixed #198 -- ensure new users have is_active == True. (@mboersma)
- [#200](https://github.com/deis/deis/pull/200) Merge pull request #200 from opdemand/198-user-isnt-active (@opdemand)
- [8210ff9](https://github.com/deis/deis/commit/8210ff9e289ef043901aa7ddd54d04240a98ea24) Fixed #202 -- adding content-type header broke registration. (@mboersma)
- [68f40b3](https://github.com/deis/deis/commit/68f40b31aca499e576316bdd1e4b49b64ded5f01) Fixed #170 -- "deis nodes:create" allows adding external instances. (@mboersma)
- [#201](https://github.com/deis/deis/pull/201) Merge pull request #201 from opdemand/170-manually-add-nodes (@opdemand)
- [0d43479](https://github.com/deis/deis/commit/0d43479b095ea1348a0b810f7de010b31560ce0f) Fixed #56 -- implemented account cancellation. (@mboersma)
- [#203](https://github.com/deis/deis/pull/203) Merge pull request #203 from opdemand/56-account-cancellation (@opdemand)
- [29a3a09](https://github.com/deis/deis/commit/29a3a093b28df33ec40d430bee8a748780d5e8d6) Fixed #205 -- updated with fresher AMIs. (@mboersma)
- [8a6f44e](https://github.com/deis/deis/commit/8a6f44e78bb25b32786ac0bb648f9ae28dbbd75d) Fixed #204 -- `deis containers` works before `git push` without error. (@mboersma)
- [#207](https://github.com/deis/deis/pull/207) Merge pull request #207 from opdemand/205-update-amis (@opdemand)
- [ef59b7c](https://github.com/deis/deis/commit/ef59b7cf36dbfc84365adc83e784702665f71a08) Fixed #122 -- CLI-driven test suite. (@mboersma)
- [88bd938](https://github.com/deis/deis/commit/88bd938e9048e6986a51606a0864fca07284f274) Updated tests (@mboersma)
- [a531dbc](https://github.com/deis/deis/commit/a531dbcd8cacde233035701a552c980bd0527a09) Merge branch '122-acceptance-test' of https://github.com/opdemand/deis into 122-acceptance-test (@mboersma)
- [d071885](https://github.com/deis/deis/commit/d0718855557961a0a30f2da75431f9861e199890) Updated version strings to 0.1.0. (@mboersma)
- [5ffb41f](https://github.com/deis/deis/commit/5ffb41f1d15ed29024abe90b9733a821eb9d630b) write knife output to celery logs regardless of success or failure fixes #206 (@gabrtv)
- [#209](https://github.com/deis/deis/pull/209) Merge pull request #209 from opdemand/log-knife-output (@opdemand)
- [b58cb38](https://github.com/deis/deis/commit/b58cb385686c45ca5c5a2f93c5986f28de96eafc) Changed DEIS_SERVER import error to a warning message. (@mboersma)
- [639c530](https://github.com/deis/deis/commit/639c53047a53a3780320821294dd3dcc8a9eee48) Remove client tests from default targets; they're too expensive. (@mboersma)
- [24d6aa7](https://github.com/deis/deis/commit/24d6aa78e3f1eec3b259f2060a2d65b0a94cd130) update docs, and remove line break from command output (@gabrtv)
- [96b7251](https://github.com/deis/deis/commit/96b72518c0ecb94fb581c86fac935a33aad4127b) update berksfile and gemfile for 0.1.0 (@gabrtv)
- [#210](https://github.com/deis/deis/pull/210) Merge pull request #210 from opdemand/122-acceptance-test (@opdemand)
- [#213](https://github.com/deis/deis/pull/213) Merge pull request #213 from opdemand/update-docs (@opdemand)
- [#212](https://github.com/deis/deis/pull/212) Merge pull request #212 from opdemand/update-ruby-deps (@opdemand)
- [0cf1279](https://github.com/deis/deis/commit/0cf127953c79032893c739b980ae2c833ab26a6d) print converge output for every node converge, success or fail (@gabrtv)
- [59706e1](https://github.com/deis/deis/commit/59706e172818230a903508af59a3553c3794abbe) remove converge from release handling, fixes double-converge issue on build/config changes (@gabrtv)
- [#216](https://github.com/deis/deis/pull/216) Merge pull request #216 from opdemand/fix-double-converge (@opdemand)
- [e829988](https://github.com/deis/deis/commit/e829988d54e10a9ceb4f4ff99531cceb60b3c6dd) Updated client README. (@mboersma)

### v0.0.7 (2013/09/10 17:06 +00:00)
- [9420dfd](https://github.com/deis/deis/commit/9420dfd017b6232e31040f5de0119a24deef619c) switch back to dev release for 0.0.7 (@gabrtv)
- [f6de3f4](https://github.com/deis/deis/commit/f6de3f43ad2c9d2445a69633e10b78ce041f0b51) Updated release docs and bumped CLI version to 0.0.7 dev. (@mboersma)
- [1274278](https://github.com/deis/deis/commit/127427894d1e52edc441b54f337434cd137f15cb) Closes #126. Updated RHS sidebar to float right. Updated the social bar to adjust its top margin on page resize. Updated typography in RHS sidebar. Updated social bar to align to bottom of page text. Updated search results page styling. (@bengrunfeld)
- [1fe584f](https://github.com/deis/deis/commit/1fe584f3d67f5000f36138b2312d9b0d7d280edc) Fixed styling error on short documentation pages (@bengrunfeld)
- [3ea1730](https://github.com/deis/deis/commit/3ea1730d749fbf92796156db58e5be44bc761cb3) Fixed #127 -- use re.search since re.match always matches string start. (@mboersma)
- [3f73d81](https://github.com/deis/deis/commit/3f73d811fa968965d053aee0d6ce696297f9398e) Removed empty docs for intentionally empty web/models.py (@mboersma)
- [f39246a](https://github.com/deis/deis/commit/f39246a3cd1f768f4c2a84b15f09e719aa218860) Fixed #128 -- @task-decorated functions need autofunction for docs (@mboersma)
- [#133](https://github.com/deis/deis/pull/133) Merge pull request #133 from opdemand/128-sphinx-task-decorator (@opdemand)
- [784358a](https://github.com/deis/deis/commit/784358aa1e4c37dca58c73e27bb1c01b4fde0918) Fixed #131 -- better error handling in `deis open` (@mboersma)
- [85afc09](https://github.com/deis/deis/commit/85afc0979279f59b3e2a9a5650b2b72b4c0f8287) Fixed #129 -- work around EC2's laggy create_security_group(). (@mboersma)
- [5150a73](https://github.com/deis/deis/commit/5150a73e77450e4bfdc5b7641697f7c79446a9f2) Fixed #135 -- better container allocation when removing nodes. (@mboersma)
- [9704d86](https://github.com/deis/deis/commit/9704d86dfa028da8cca719ea5745e152c2761100) Added container allocation tests, re #135. (@mboersma)
- [5de5cd1](https://github.com/deis/deis/commit/5de5cd1280199342b6cd80c50430440ec7401688) Fixed stupid flake8 errors. (@mboersma)
- [3ae5b4c](https://github.com/deis/deis/commit/3ae5b4c5fbe567f5dad1ae7652c8d8a0f16267e1) Fixed #140 -- more compatible usage of mktemp in provision script. (@mboersma)
- [259134c](https://github.com/deis/deis/commit/259134c9e136b29f2fdf6902007512f1800e79c0) Fixed searchbar breaking when user inputs too much text (@bengrunfeld)
- [c2fccb7](https://github.com/deis/deis/commit/c2fccb7f0a56ea470f92a43f6b2f830053e96e50) Fixes #141. Adds sitemap to docs. jQuery for menu functionality adjusted to speed up load time. (@bengrunfeld)
- [c2f9f28](https://github.com/deis/deis/commit/c2f9f28239c5f9a6204b2931074a94aff974b100) Updated comments on adjustment.js file (@bengrunfeld)
- [2c314c9](https://github.com/deis/deis/commit/2c314c912f49253bf458a87672b84b2b9ab71fa2) Fixes #142. Meta tag containing noindex,nofollow removed. (@bengrunfeld)
- [ed7e044](https://github.com/deis/deis/commit/ed7e0440806c4260b10ced94180fff86dae6a3a0) Updated gunicorn and django-json-field. (@mboersma)
- [e4110c3](https://github.com/deis/deis/commit/e4110c3794ddf2316a6b0fa9916ffc8d4d09203d) Fixes #143. Adjusted canoncial tag to point to correct URL (@bengrunfeld)
- [44a5a9e](https://github.com/deis/deis/commit/44a5a9e53d84ae955c10838564241a1c47c7ea6d) Fixed #145 -- allow empty comment/email in SSH key regex. (@mboersma)
- [9ba636d](https://github.com/deis/deis/commit/9ba636dbc4a5b93dd4231cf3ad0b41facbcd9a24) Fixed #148 -- stick with celery 3.0.22 for now. (@mboersma)
- [#149](https://github.com/deis/deis/pull/149) Merge pull request #149 from opdemand/148-new-celery-breaks (@opdemand)
- [#146](https://github.com/deis/deis/pull/146) Merge pull request #146 from opdemand/145-ssh-key-regex (@opdemand)
- [e40d0aa](https://github.com/deis/deis/commit/e40d0aa2371012d0218e33a139f080132bbc0b8d) Fixed #153 -- warn and return if no SSH keys found. (@mboersma)
- [#154](https://github.com/deis/deis/pull/154) Merge pull request #154 from opdemand/153-no-ssh-keys (@opdemand)
- [32d61e5](https://github.com/deis/deis/commit/32d61e54793d0f164ef45e44cae9b9b054513051) add timeout/attempts to util.connect_ssh, use 120s by default (@gabrtv)
- [1c39a12](https://github.com/deis/deis/commit/1c39a12c536f666f589aed28ac120e63ca9c1e71) Added the develop branch to Travis CI. (@mboersma)
- [#155](https://github.com/deis/deis/pull/155) Merge pull request #155 from opdemand/fix-ec2-ssh-timeout (@opdemand)
- [86227fa](https://github.com/deis/deis/commit/86227fa39e9e5153eebec0348f53aad761c33f2d) raise EnvironmentError if no credentials provided (@gabrtv)
- [0a47d7c](https://github.com/deis/deis/commit/0a47d7cef83b887a14e71ecdc40dc1ac3cf032f5) raise 400 on environment error, as in case of missing credentials (@gabrtv)
- [bea8963](https://github.com/deis/deis/commit/bea89638217c03649493701e9b67c83f8df7ae2c) handle no credentials error on formation/node destroy (@gabrtv)
- [#158](https://github.com/deis/deis/pull/158) Merge pull request #158 from opdemand/fix-missing-credentials (@opdemand)
- [dc7ae67](https://github.com/deis/deis/commit/dc7ae672a449f526336602c5f7224d9597692b66) Fixed #124 -- fill in explicit values for Flavor params. (@mboersma)
- [#159](https://github.com/deis/deis/pull/159) Merge pull request #159 from opdemand/124-ami-in-flavors (@opdemand)
- [b466011](https://github.com/deis/deis/commit/b466011d738f1a4335b906c1678d7199320d6762) add uniqueness constraint on public key field, with migration and test (@gabrtv)
- [#160](https://github.com/deis/deis/pull/160) Merge pull request #160 from opdemand/fix-duplicate-sshkey (@opdemand)
- [1e0cbe6](https://github.com/deis/deis/commit/1e0cbe622973875cc7ffaa0e87e03717b63046a3) add 500.html template as part of api staticfiles (@gabrtv)
- [#161](https://github.com/deis/deis/pull/161) Merge pull request #161 from opdemand/add-500-template (@opdemand)
- [b209957](https://github.com/deis/deis/commit/b20995746ca5117229aad45bbd9e2995651f7374) Fixed #144 -- check provisioning dependencies before running script (@mboersma)
- [#162](https://github.com/deis/deis/pull/162) Merge pull request #162 from opdemand/144-check-deis-deps (@opdemand)
- [9433bf2](https://github.com/deis/deis/commit/9433bf2efc903f0402651be135d3dd480a8b8d43) retry all chef api calls with configurable attempts and interval (@gabrtv)
- [#163](https://github.com/deis/deis/pull/163) Merge pull request #163 from opdemand/feature/chef-api-retries (@opdemand)
- [fe70977](https://github.com/deis/deis/commit/fe70977ea2fc46283c9c3e1a6baf320aa49f8b43) Fixed #123 -- updated EC2 AMIs and pointers to them. (@mboersma)
- [#164](https://github.com/deis/deis/pull/164) Merge pull request #164 from opdemand/123-update-amis (@opdemand)
- [1614736](https://github.com/deis/deis/commit/16147368a222c4a756f5e6733a7bcd51f20d4ed2) Added tests for Flavor updating. (@mboersma)
- [7fc1413](https://github.com/deis/deis/commit/7fc141346ec01a8c38eb91c23200369a7e9326d6) Fixed #151 -- `deis flavors:update` working with tests. (@mboersma)
- [#165](https://github.com/deis/deis/pull/165) Merge pull request #165 from opdemand/151-update-flavors (@opdemand)
- [fe0283e](https://github.com/deis/deis/commit/fe0283e8a3cb42890ec0dc028cb36f4732267232) Removed "develop" from Travis CI, we're not git-flowing now. (@mboersma)
- [#168](https://github.com/deis/deis/pull/168) Merge pull request #168 from opdemand/develop (@opdemand)

### v0.0.6 (2013/08/21 15:49 +00:00)
- [3a48716](https://github.com/deis/deis/commit/3a48716595c6c33696cc03efac5004b2d4e6483e) Fixed #71 -- help sphinx work on Windows. (@mboersma)
- [e48b264](https://github.com/deis/deis/commit/e48b264a61a726027758f0af324ac1c4cbac4a28) Updated jQuery page height detection. Updated styling: added rollover functionality. #29 (@bengrunfeld)
- [dcc2859](https://github.com/deis/deis/commit/dcc2859f7bdfb76080140cef847f1006052764f4) Updated Releases link in Docs #29 (@bengrunfeld)
- [1788290](https://github.com/deis/deis/commit/1788290a24718e350cbd675c564eca278321a481) Updated styling on docs. #29 (@bengrunfeld)
- [7cd482d](https://github.com/deis/deis/commit/7cd482d60df09a69150ab337c30ce8a36f6ddfb5) Silenced some sphinx warnings. (@mboersma)
- [7d3549c](https://github.com/deis/deis/commit/7d3549c8e4a409078e36c852b818817b470b9ada) Updated releases doc. (@mboersma)
- [c3d1800](https://github.com/deis/deis/commit/c3d1800f7c991bf6b6795eaa407e3fb4103f0c99) Add meta tags for Terms section #61 (@JoshuaSchnell)
- [be982a9](https://github.com/deis/deis/commit/be982a97c65c6812293971909f99ca21f6ef7651) Add meta tags to client reference, community, etc. #61 (@JoshuaSchnell)
- [6bd1d52](https://github.com/deis/deis/commit/6bd1d52a928f9c3d2b812d83cda5ceefb7d30d9e) Add meta to getting started and index #61 (@JoshuaSchnell)
- [882940b](https://github.com/deis/deis/commit/882940be8113925cfbc267bfd7a3f3e653aba6eb) Fixed #59 -- restore boilerplate JS to fix Sphinx quick search. (@mboersma)
- [671e0de](https://github.com/deis/deis/commit/671e0de7467584bfb1253c813d3bd4200c79ad87) Added local copy of searchtools.js for tweaking. (@mboersma)
- [3953b9f](https://github.com/deis/deis/commit/3953b9f44c26d9761df37ef9d4a890c1544a64da) Updated Sphinx documentation templates (@bengrunfeld)
- [ba7d671](https://github.com/deis/deis/commit/ba7d6710d530dfbab69c67b18834efa8423b8f1d) Updated documentatipon theme. Removed elements from Search Results page and adjusted the Search box in the sidebar. Changed searchtools.js to resize footer to align with bottom of content (@bengrunfeld)
- [8c2ca4d](https://github.com/deis/deis/commit/8c2ca4d9deb5b6660de1b9d8c8cb7156095e6039) Fixed #67 -- synced up README.rst with other intro docs. (@mboersma)
- [#81](https://github.com/deis/deis/pull/81) Merge pull request #81 from opdemand/67-cli-readme (@opdemand)
- [89ee756](https://github.com/deis/deis/commit/89ee75620e01387cf52d06e5ecea2a86dbe0c805) Moved checkURL.js out of the docs _build directory. (@mboersma)
- [552b4f2](https://github.com/deis/deis/commit/552b4f28c78e832123ac517532ad3770ee384c32) Fixed #61 -- finished adding META desc & keywords to server API docs. (@mboersma)
- [4edc8f6](https://github.com/deis/deis/commit/4edc8f6d4f5b4155ec3c1eb7664b5d2877310a62) Fixed #83 - use vert bar as docs title separator. (@mboersma)
- [acd23f3](https://github.com/deis/deis/commit/acd23f3c97864c52c159c295931d1b6d12233903) Removed reference to "categories" when search returns no results. (@mboersma)
- [bddf4ca](https://github.com/deis/deis/commit/bddf4caa89a79ac5d14143ce4bcf6783643cdb18) Removed extra copy of checkURL.js under sphinx theme dir. (@mboersma)
- [392c9fb](https://github.com/deis/deis/commit/392c9fb5274a7ef78d6818546cc5654c634cca18) Change example IP address per RFC 5737 (@mboersma)
- [de4f700](https://github.com/deis/deis/commit/de4f700814a024fe2271863651f48a5503d1871e) Updated South to 0.8.2. (@mboersma)
- [594d68c](https://github.com/deis/deis/commit/594d68c9eca47063f5738728f58088b63a4b43a4) Fixed #66 -- developer setup docs done. (@mboersma)
- [5e0f219](https://github.com/deis/deis/commit/5e0f21959bb1ef71301ef66afad6a2c5b732fda9) Updated to Django 1.5.2 security release, also boto 2.10.0. (@mboersma)
- [f6a46f1](https://github.com/deis/deis/commit/f6a46f1908c0ad61f5a2adb2bdfd8d478e1e5023) Deis does not rest. (@mboersma)
- [#85](https://github.com/deis/deis/pull/85) Merge pull request #85 from opdemand/66-dev-setup (@opdemand)
- [9898dd8](https://github.com/deis/deis/commit/9898dd853ad9536c52339da5ebb66efeec07a7a8) Fixed #38 -- enabled Django admin for Deis API models. (@mboersma)
- [#90](https://github.com/deis/deis/pull/90) Merge pull request #90 from opdemand/38-django-admin (@opdemand)
- [371638b](https://github.com/deis/deis/commit/371638b51e33c9d4112db030a27743356608d434) Fixed #92 -- correct some JavaScript in docs for "Releases" link. (@mboersma)
- [14481e3](https://github.com/deis/deis/commit/14481e3f4793103c89de288f99e9f91c871cbc49) Fixed #93 -- web UI loads the intended assets now. (@mboersma)
- [97c6030](https://github.com/deis/deis/commit/97c60305afe9800af2cb17e929523aa813ee0eb3) Fixed #94 -- reference squashing commits in docs, many thanks to Docker for the general wording. (@mboersma)
- [34cfcfd](https://github.com/deis/deis/commit/34cfcfdbd43e16d4cc6cd00ef17aee6b93f00460) Fixed 2 broken URL refs, re #94. (@mboersma)
- [e1ad061](https://github.com/deis/deis/commit/e1ad06143b991c563d323df6652deaafcd59d0f5) use deis::controller meta recipe in provision script re #95 (@gabrtv)
- [40c1676](https://github.com/deis/deis/commit/40c1676be588a9cbc965cf1af69aeb872e167d7b) Fixed #80 -- serializer SlugRelatedFields used too broad a query. (@mboersma)
- [#97](https://github.com/deis/deis/pull/97) Merge pull request #97 from opdemand/80-slug-related-field (@opdemand)
- [99729b1](https://github.com/deis/deis/commit/99729b1224c5b871876909d905921f42443478f3) add logs endpoint, views and tests along with `deis logs` cli command re #20 (@gabrtv)
- [#98](https://github.com/deis/deis/pull/98) Merge pull request #98 from opdemand/log-aggregation (@opdemand)
- [9b9ca9a](https://github.com/deis/deis/commit/9b9ca9a25eec9e4cb8e4429c27f74010bc85d0f2) Updated django-restframework to 2.3.7. (@mboersma)
- [aae91eb](https://github.com/deis/deis/commit/aae91ebb2ece10e4f4e292976e897bc769a5564c) Expanded vocabulary. (@mboersma)
- [ce8d1be](https://github.com/deis/deis/commit/ce8d1be44910fa03a1991202f360b9fab24cd159) Fixed #100 -- added rsyslog port 514 to EC2 provisioning script. (@mboersma)
- [06c4347](https://github.com/deis/deis/commit/06c4347775f93f5ea99d9a6c3393f100cf2158a8) Fixed #96 -- doc the 2 restframework methods. (@mboersma)
- [b6cf21f](https://github.com/deis/deis/commit/b6cf21f5d24df17021ff1eceed1a18f295e10e9a) Fixed #78 -- cookie cleanup before login. (@mboersma)
- [6fab360](https://github.com/deis/deis/commit/6fab3601f08d849b860d6615cccd73b639a5891c) Fixed #68 -- better CLI handling of EnvironmentErrors. (@mboersma)
- [#102](https://github.com/deis/deis/pull/102) Merge pull request #102 from opdemand/78-login-failed (@opdemand)
- [#103](https://github.com/deis/deis/pull/103) Merge pull request #103 from opdemand/68-cli-stacktrace (@opdemand)
- [1ec6ae4](https://github.com/deis/deis/commit/1ec6ae4df38b55aed17d8b3e80a096e17339526a) Added `deis logs` to REST API document, re #20. (@mboersma)
- [7132af7](https://github.com/deis/deis/commit/7132af76016391154aa8d6e7524d95b459dd76fc) switch from dash to bash, use /usr/bin/env for better compat fixes #104 (@gabrtv)
- [f392868](https://github.com/deis/deis/commit/f392868170d017df23a02b92c8f38fab0c64af68) switch berksfile to development version (@gabrtv)
- [3865c57](https://github.com/deis/deis/commit/3865c57ff4c759cd65d9446fcfad5ae97bb75742) dont converge on any job, just on new node creation fixes #99 (@gabrtv)
- [#105](https://github.com/deis/deis/pull/105) Merge pull request #105 from opdemand/99-autherror-on-scale (@opdemand)
- [69526b0](https://github.com/deis/deis/commit/69526b0f956a117cd8c87566719a93ca8fde33ef) add logs to common CLI commands (@gabrtv)
- [e71c803](https://github.com/deis/deis/commit/e71c803422cf77e3678ceae0b7d5b8d3ab8177dd) find a proxy and use the OS default handler to open the URL fixes #106 (@gabrtv)
- [#107](https://github.com/deis/deis/pull/107) Merge pull request #107 from opdemand/106-deis-open (@opdemand)
- [30a35e8](https://github.com/deis/deis/commit/30a35e8adb1adf071db8feef356d94aa09b4e762) log django.request and api errors to console, where gunicorn will dump a stacktrace fixes #79 (@gabrtv)
- [f26e859](https://github.com/deis/deis/commit/f26e85912a7b62a326d834a0bcf41cdc0f4a3154) return a proper 404 when no logs exist fixes #101 (@gabrtv)
- [#109](https://github.com/deis/deis/pull/109) Merge pull request #109 from opdemand/101-no-logs-message (@opdemand)
- [6081551](https://github.com/deis/deis/commit/608155165fefa13b2bf5e0c07efb5a35273a68aa) cleanup any old logfiles re #101 (@gabrtv)
- [ca3112c](https://github.com/deis/deis/commit/ca3112c0cb42991fb00d423cb9838a7486335593) add whitespace after : for the flake8 gods (@gabrtv)
- [#110](https://github.com/deis/deis/pull/110) Merge pull request #110 from opdemand/101-no-logs-message (@opdemand)
- [#108](https://github.com/deis/deis/pull/108) Merge pull request #108 from opdemand/79-log-500-errors (@opdemand)
- [30a939f](https://github.com/deis/deis/commit/30a939f8132ef5760378460ae8e7803aff0ee753) Added `deis logs` to docstrings. (@mboersma)
- [313200c](https://github.com/deis/deis/commit/313200c8abfd2f46a51ddf20f589dd2004925d52) remove dupe doc reference to logs (@gabrtv)
- [2d2c65d](https://github.com/deis/deis/commit/2d2c65d22d6a9903a84c823d9e32059dbbc61ec8) Fixed #84 -- no stacktrace when `deis providers` comes up empty. (@mboersma)
- [e26da5e](https://github.com/deis/deis/commit/e26da5e5b70244243adffd8c05a41ef5d2d963df) Fixed #24 -- added progress animation for long CLI commands. (@mboersma)
- [#112](https://github.com/deis/deis/pull/112) Merge pull request #112 from opdemand/24-cli-progress (@opdemand)
- [976ba88](https://github.com/deis/deis/commit/976ba8849331bafd6663e718d5529dfbca68b6d7) Aligned the Deis web UI with the deis.io/docus theme. (@mboersma)
- [#114](https://github.com/deis/deis/pull/114) Merge pull request #114 from opdemand/113-controller-web-ui (@opdemand)
- [bf61fbc](https://github.com/deis/deis/commit/bf61fbc03f2edf3ef9ef3eb595f04b05e0079c72) Fixed #111 -- move new_nodes out of conditional and delete nodes properly. (@mboersma)
- [b330e24](https://github.com/deis/deis/commit/b330e2422811c1d2f011742e41837357cb8fc032) Fix #70 -- first pass at `deis run <command>` including test coverage (@mboersma)
- [#116](https://github.com/deis/deis/pull/116) Merge pull request #116 from opdemand/111-scaling-error (@opdemand)
- [#115](https://github.com/deis/deis/pull/115) Merge pull request #115 from opdemand/70-run-bash (@opdemand)
- [4033d94](https://github.com/deis/deis/commit/4033d9406c85f749453a65e702222cefbab9315e) Fixed #118 -- fixed progress thread management in CLI. (@mboersma)
- [#119](https://github.com/deis/deis/pull/119) Merge pull request #119 from opdemand/118-progress-threads (@opdemand)
- [ce2563b](https://github.com/deis/deis/commit/ce2563b90d3b7d4becc0cd4750602b0d19386ecb) import profile.d/*.sh into shell environment on `deis run` (@gabrtv)
- [9d7de9b](https://github.com/deis/deis/commit/9d7de9b15f2b09668686556ae5542f487b7e772f) Prep for 0.0.6 release. (@mboersma)
- [a3ab116](https://github.com/deis/deis/commit/a3ab1163cf64e21c52c0630fabec5ce6ee8eb587) update berksfile.lock for 0.0.6 (@gabrtv)

### v0.0.5 (2013/08/06 15:59 +00:00)
- [e5b7c68](https://github.com/deis/deis/commit/e5b7c686cced8fd4011e3163d22b1327326d65cc) Fixed #5 -- created release process doc. (@mboersma)
- [0ac52a8](https://github.com/deis/deis/commit/0ac52a84a0a5d4c33fa9ba5628f5856ffd18787d) Add client module to Sphinx docs, re #49. (@mboersma)
- [b455a1d](https://github.com/deis/deis/commit/b455a1dbf1f04823228990b5870ecc5cbf7cacad) reset south migrations, and add uniqueness contraint on builds fixes #30 (@gabrtv)
- [ab88961](https://github.com/deis/deis/commit/ab8896176ddcc3c988594baed06dd04fff9c522e) paginate by 100 instead of 10 fixes #44 (@gabrtv)
- [f03d218](https://github.com/deis/deis/commit/f03d218ee5e3e34cedb4d4b512b169ce57f8831c) display assigned nodes on containers:list, to view balancing (@gabrtv)
- [849ebc5](https://github.com/deis/deis/commit/849ebc5fd93c6718a23e91f0aa0d93cb54ea8769) Added brief REST API docs, re #47. (@mboersma)
- [56effef](https://github.com/deis/deis/commit/56effefb8f763b090417cf82a5d73a67fd2eb698) Fixed URL ordering that broke several tests. (@mboersma)
- [df32c09](https://github.com/deis/deis/commit/df32c090c2fcbcb269514c118c38e0123f023169) Fixed #49 -- make module docstring into reST, everyone's happy. (@mboersma)
- [213237f](https://github.com/deis/deis/commit/213237fc4fe775e549bb78c64b0cfac5d351873b) Simplify requirements and fake SECRET_KEY for doc generation, re #57. (@mboersma)
- [b4dc254](https://github.com/deis/deis/commit/b4dc2549856f1bb97f31aeb4998f77e25593f026) Include LICENSE file in client pypi package. (@mboersma)
- [5faf53f](https://github.com/deis/deis/commit/5faf53f29c90e745aa4390c681db02eca4db2719) Updated sphinx theme to match deis.io site, re #29. (@mboersma)
- [51434ae](https://github.com/deis/deis/commit/51434ae1d1a7779545eb9297b3ac9b336da9c155) Added right navigation to docs. (@mboersma)
- [f5891f2](https://github.com/deis/deis/commit/f5891f2c387c70715736512c8dc4f5aa4ced3578) wait 10 seconds for ssh daemon to come up (@gabrtv)
- [9da39f7](https://github.com/deis/deis/commit/9da39f7ed60d5fee31a73fecd447d6dfde99ecdb) Added sphinx search in for styling, re #29. (@mboersma)
- [88c98c6](https://github.com/deis/deis/commit/88c98c61e53e8c9b8ba59087836535688e20a03e) Updated: CSS and Structure (@bengrunfeld)
- [335b0ad](https://github.com/deis/deis/commit/335b0ad0e74ba0024818aafb61919d4daa85237d) Deleted old CSS files. Updated Layout and CSS (@bengrunfeld)
- [468ca52](https://github.com/deis/deis/commit/468ca52f925920c63a99f341b04dc3780a18fb11) Remove permalinks from Sphinx doc generation. (@mboersma)
- [cc96896](https://github.com/deis/deis/commit/cc968961acbd2bd00a2c52d350861cc4fee8472a) Fixed broken JS by removing a macro. (@mboersma)
- [635988a](https://github.com/deis/deis/commit/635988af2cdeb51c4d6054798fec3f901bcfe32e) Updated: Margins and Padding inside of Sidebar (@bengrunfeld)
- [12d4e13](https://github.com/deis/deis/commit/12d4e13c831728a25f9006efc6ac3dd7ad444947) Merge remote-tracking branch 'origin/master' (@bengrunfeld)
- [f5be86c](https://github.com/deis/deis/commit/f5be86cf9345d703fdfde5ea2bd0715887414d4e) Restored search in docs, started Community docs, re #4. (@mboersma)
- [6f11be6](https://github.com/deis/deis/commit/6f11be60d583ebb254ccb45fd84356362d94265a) Added: Analytics code to layout.html (@bengrunfeld)
- [d518536](https://github.com/deis/deis/commit/d518536feb3730b48bf63874cd066694c1714627) Added a "make zipfile" target for pythonhosted.org doc hosting. (@mboersma)
- [44b3c7d](https://github.com/deis/deis/commit/44b3c7d88bf995d91862e33f2eef188092701fe4) Update releases doc to include "make zipfile" and pythonhosted.org. (@mboersma)
- [652af70](https://github.com/deis/deis/commit/652af70fdf9e7a164da9bba310b22d206c432f12) Updated: Removed search and restyled sidebar. Code blocks now in lightgray background (@bengrunfeld)
- [013ea3a](https://github.com/deis/deis/commit/013ea3a43f5a87d8b4efbe2677fc9efeae6718f4) Fixed: Javascript error. Changed docs styles (@bengrunfeld)
- [e0db967](https://github.com/deis/deis/commit/e0db967955a27aaaf9225f934ff448a32bf0d96b) Fixed #4 -- added community and conduct language to docs. (@mboersma)
- [639e80c](https://github.com/deis/deis/commit/639e80cf21a6771539e9d10e085e3d62a5143ed1) pushing new structure for sphinx nav (@gabrtv)
- [7fd6f14](https://github.com/deis/deis/commit/7fd6f141a5e2ae3989686a0b1508aac24cba9cc0) add comprehensive client reference, cleanup toctrees re #47 (@gabrtv)
- [cac60eb](https://github.com/deis/deis/commit/cac60eb4431e30a9745f4fbc4bcfdd2238cb6856) first pass at technical overview, include terms in toctree #62 (@gabrtv)
- [93d45fc](https://github.com/deis/deis/commit/93d45fcdd544fa68e03e4e3e53ed9779e6847205) Updated: the jQuery that sets the height of the 3 columns, and how the margin above teh social buttons is calculated (@bengrunfeld)
- [d47b85d](https://github.com/deis/deis/commit/d47b85de5b55f1d491022ecc8bdda3188b76179b) Added accordion functionality to docs sidebar. Updated style of search box. Updated jQuery page height calculation. (@bengrunfeld)
- [3f066ee](https://github.com/deis/deis/commit/3f066eef02668ff36b847a87e0cbdea62e42f1ed) fix seo on main documentation page (@gabrtv)
- [4caf32c](https://github.com/deis/deis/commit/4caf32cd10c26a09dcc1bcadee8b7e9a2f5ac079) rename technical overview to concepts (@gabrtv)
- [bec19a0](https://github.com/deis/deis/commit/bec19a0bc88b20acd4b18254c56204fa7b874ac3) Added the search results template to deis theme. (@mboersma)
- [65d8b99](https://github.com/deis/deis/commit/65d8b9922de74583fc5bfd5e126d86cce377a77a) Started dev setup and code standards docs, re #2. (@mboersma)
- [25e9cc2](https://github.com/deis/deis/commit/25e9cc25fbf10dbddb3e22738bbac75a893d752d) update concepts language, welcome page language (@gabrtv)
- [9eed2bf](https://github.com/deis/deis/commit/9eed2bf4a30baf81e80d142c05f9283b988f493c) Fixed #2 -- added "contributing" documentation. (@mboersma)
- [ae296bd](https://github.com/deis/deis/commit/ae296bd3a595e26a33a8b1732cf5f552077e13e0) standardize on uppercase for contributing docs (@gabrtv)
- [e0270a2](https://github.com/deis/deis/commit/e0270a2e2c03b9b01e4f8adf8de8db782d4ed9de) fix some typos and broken links on concepts page (@gabrtv)
- [6d8f971](https://github.com/deis/deis/commit/6d8f971a67a1f97cc8dfad6c38a964ec570e8816) add technical description to welcome page (@gabrtv)
- [117b6d2](https://github.com/deis/deis/commit/117b6d2d0d182eca5d8a6d7a0ee4d4822be8e765) update client readme.rst with updated language (@gabrtv)
- [8743735](https://github.com/deis/deis/commit/8743735aec9f955231c0e99153028cfb2af1247f) add getting started instructions, remove CRs fixes #45 (@gabrtv)
- [e7ab6d6](https://github.com/deis/deis/commit/e7ab6d6ba6cc231abef2423dc630541d862efcd4) Started dev setup documentation, re #1. (@mboersma)
- [7c43037](https://github.com/deis/deis/commit/7c43037c181882471db6e35ea3814df984523ef6) add better client description (@gabrtv)
- [a1cc87c](https://github.com/deis/deis/commit/a1cc87c6293741bfd77268467337d0f790092231) edits to readme (@gabrtv)
- [08bc8e1](https://github.com/deis/deis/commit/08bc8e17ab8135ee7fb6d0e8474030666d16ecb0) Updated README.md with some github-flavored markdown (@mboersma)
- [b778f81](https://github.com/deis/deis/commit/b778f81b4ad222667f3705f3e1510e64f2cab70d) Updated jQuery URL-matching functionality to sidebar (@bengrunfeld)
- [cbe5145](https://github.com/deis/deis/commit/cbe51452ecce807aa1fdb56707061d699a61ae29) Updated docs layout. Removed old code. (@bengrunfeld)
- [a4ffa1f](https://github.com/deis/deis/commit/a4ffa1f373bd78f8da0cfe234108bd4143ca1cda) Added the installation doc, from the project README. (@mboersma)
- [d29c97a](https://github.com/deis/deis/commit/d29c97a57fe3ca73ee08db8e4f614867f5106382) Add ref to Installation from Usage doc. (@mboersma)
- [e91f40d](https://github.com/deis/deis/commit/e91f40db9adee8dca5e0189d84c4c485537b3791) Created new favicon and applied it to docs and theme (@bengrunfeld)
- [5b4f5bc](https://github.com/deis/deis/commit/5b4f5bcf589601003b39cebc5e4855f7256e8a0d) Reference the Installation doc in Developer Setup. (@mboersma)
- [a6a4aa1](https://github.com/deis/deis/commit/a6a4aa11909b6370a9eecef30202c2fbf683321a) Fixed link to example app and fixed some inline <pre> markup. (@mboersma)
- [5df9616](https://github.com/deis/deis/commit/5df96161edd366589209cba4f23b8515c90e5b6c) add note about publishing amis (@gabrtv)
- [a39574f](https://github.com/deis/deis/commit/a39574f82a0c8eee044d54cc32a52d8e32c65549) first pass at terms documentation (@gabrtv)
- [5b72fe8](https://github.com/deis/deis/commit/5b72fe879c4e413fdf24b3b8b2edfba8bcb6afce) Added module docstrings for the proto-web app. (@mboersma)
- [29a751c](https://github.com/deis/deis/commit/29a751c3326ae310a15801c88aa6be8fd2c68950) Skip 0-length files and use Deis' CSS in coverage report. (@mboersma)
- [8d4fba5](https://github.com/deis/deis/commit/8d4fba567de967b31e2498146ae5eff3ca71d89b) Commented out empty web test. (@mboersma)
- [0d8474a](https://github.com/deis/deis/commit/0d8474a43ac6086795378081b402a3c3d9e1fb97) Updated boto, django-celery, and paramiko. (@mboersma)
- [42bf7ef](https://github.com/deis/deis/commit/42bf7ef6ff299f7a551bcb447166d4e4b368479f) remove dupe content from docs welcome page (@gabrtv)
- [f9107e3](https://github.com/deis/deis/commit/f9107e32a3eb316f9257c0b630f8f897abde38d4) add note about admins group requirement re #53 (@gabrtv)
- [bd67434](https://github.com/deis/deis/commit/bd6743448e0395b1f2ecc046f05abce02942ba5d) update installation docs w/ chef admin requirement re #53 (@gabrtv)
- [285588c](https://github.com/deis/deis/commit/285588ce626aaea57bf3572bc0ab3de5e075c298) Updated version to 0.0.5. (@mboersma)
- [266fa79](https://github.com/deis/deis/commit/266fa79dfb6726cd9c48aa89c0521167a456c2a4) Removed a swapfile dropping and updated docs Makefile. (@mboersma)
- [868d037](https://github.com/deis/deis/commit/868d0378c29a713e9a29f31275088b96ee6d5437) update Berksfile to use new v0.0.5 cookbook (@gabrtv)

### v0.0.4 (2013/07/30 15:32 +00:00)
- [0462cef](https://github.com/deis/deis/commit/0462cef5812ce31fe12f25596ff68dc614c708af) Initial commit (@mboersma)
- [40bf6f6](https://github.com/deis/deis/commit/40bf6f64ddb2d9dbe7acbd4fd528c8087d9bbca0) Merged three projects into one. (@mboersma)
- [918a18f](https://github.com/deis/deis/commit/918a18f1da8aa834ab0930cd46b572e11b25864a) Updated README.md for current deis product blurb. (@mboersma)
- [3ac0199](https://github.com/deis/deis/commit/3ac0199a689869cc34e63d5581ac40b4d5a03968) Add missing bin/ dirs and remove bin from .gitignore. (@mboersma)
- [9413a72](https://github.com/deis/deis/commit/9413a721d4ba583ef5fe2086a18d41d2567d32d6) Updated git checkout path for controller. (@mboersma)
- [1976aa6](https://github.com/deis/deis/commit/1976aa61c38c133ea4b74cc7c5fcd8616cf95ea1) More project structure maneuvering. (@mboersma)
- [d9285ce](https://github.com/deis/deis/commit/d9285ceb14b2cf1c6a76d4aacc1619c8a7de551b) change contrib script name (@gabrtv)
- [98b14f0](https://github.com/deis/deis/commit/98b14f0209431177e3505eb0c853e5e207d84c16) add Gemfile and Berksfile for Chef dependency management (@gabrtv)
- [2739380](https://github.com/deis/deis/commit/27393800efe41283cc7baa6a7c1d4a3a8aa6d4b8) add helper script for provisioning an ec2 controller (@gabrtv)
- [9b321f8](https://github.com/deis/deis/commit/9b321f8993e4318c780f7f0fb76acd7db93c3b03) ignore rbenv version (@gabrtv)
- [aa2b715](https://github.com/deis/deis/commit/aa2b715b6ef6ec9aef6afa2a021cf90ecf581ad4) remove pydevd from master (@gabrtv)
- [60814b9](https://github.com/deis/deis/commit/60814b90c68775d136ab82b56abc488791fa696b) Updated doc API layout. (@mboersma)
- [75f2aa2](https://github.com/deis/deis/commit/75f2aa20bb5e1bbf18876de3cf7a7a83964b6572) move auto-generated ssh keys to formation, add tests for ssh key override on formation create.  fixes #17 (@gabrtv)
- [7064132](https://github.com/deis/deis/commit/70641325c66dacc29cefdde0059b7a8da2561eea) Added config file for Travis CI. (@mboersma)
- [7392f83](https://github.com/deis/deis/commit/7392f83ef77c311841b1620c98d3750a2dffc763) hardcode chef version to 11.4.4 as the new 11.6.0 breaks everything (@gabrtv)
- [35c3046](https://github.com/deis/deis/commit/35c304673a91111d7aacad4ea3829bbf1d793491) update berksfile to use deis 0.0.4 (@gabrtv)
- [b58281b](https://github.com/deis/deis/commit/b58281bfd3882328a6c1c14386e97246f1016f70) More testing of Travis CI configuration. (@mboersma)
- [d7c60dc](https://github.com/deis/deis/commit/d7c60dcc9204301696448ed475f73ff6c27e3662) Added DB config to .travis.yml. (@mboersma)
- [d03aff3](https://github.com/deis/deis/commit/d03aff38ab358ea65dc2cc093d1df2886deda091) Fixed #9 -- Travis-CI.org integration working. (@mboersma)
- [9328e5a](https://github.com/deis/deis/commit/9328e5abf045950c76f6995d09fc92c2c697dd5e) Restore lost line in travis config. (@mboersma)
- [dc62bf3](https://github.com/deis/deis/commit/dc62bf3b1d916e2490d9d396f84b831d0758071e) Fixed #7 -- updated README and published to pypi. (@mboersma)
- [71bd371](https://github.com/deis/deis/commit/71bd371ed16437b2343cfbbe94902c995d8dd5e2) Updated pypi package status to beta. (@mboersma)
- [4ab89bb](https://github.com/deis/deis/commit/4ab89bbc69b97177542569f1db347245e8090cba) convert deis to symlink (@gabrtv)
- [cc9f277](https://github.com/deis/deis/commit/cc9f277e0a6c80e814fa09545c06c17c6559bada) refactor CLI http dispatch re #16 (@gabrtv)
- [485df36](https://github.com/deis/deis/commit/485df366f62979d24ae109f457fd23e2b2918158) seed default providers and flavors on user registration, re #23 (@gabrtv)
- [50df654](https://github.com/deis/deis/commit/50df65488d5df45bf70fd57edb8c03ee67e55380) add `deis providers:discover` with a more explicit discovery workflow that uses default providers, fixes #23 (@gabrtv)
- [1536176](https://github.com/deis/deis/commit/1536176f0a0aaec8aa8b38ce1cad675026652ed2) Began some PEP8 / pyflakes-inspired code cleanup. (@mboersma)
- [96fecae](https://github.com/deis/deis/commit/96fecae7962241904140fabfea1bf3136e459b3f) Continued PEP8 code cleanup. (@mboersma)
- [8b3c218](https://github.com/deis/deis/commit/8b3c218c66e512f85c27ef9fd7a3aa479a96fecd) Removed two requirements that weren't actually required. (@mboersma)
- [d97714b](https://github.com/deis/deis/commit/d97714be7632e5eaeaf51ae61ff7e842c037782e) Testing pre-commit flake8 hook. (@mboersma)
- [3b9d90a](https://github.com/deis/deis/commit/3b9d90ad024299fab31aafcd8a8c561b2bd0f2ef) Replaced pep8 and pyflakes targets with flake8. (@mboersma)
- [ba5119b](https://github.com/deis/deis/commit/ba5119be6962948d86d93a806f40cff7bf1923ec) Refs #6, flake8 code cleanup nearly done. (@mboersma)
- [0c89448](https://github.com/deis/deis/commit/0c894485bdfde4516e9950d17273c95d185b3d93) first pass at layer refactoring, with passing test suite (@gabrtv)
- [1206fe6](https://github.com/deis/deis/commit/1206fe67850c1ceaa196e2279d0721fa9f1581c5) remove south for now (@gabrtv)
- [14aa7e4](https://github.com/deis/deis/commit/14aa7e4debdc578c072279a01a6f2566e9e67942) use layer to define run_list and initial attributes (@gabrtv)
- [30e40b9](https://github.com/deis/deis/commit/30e40b9fdade14a1827c1aade514c57012f3c3a2) Refs #10 -- set up travis-ci.org and coveralls.io integration. (@mboersma)
- [c516468](https://github.com/deis/deis/commit/c516468704d381bd880d00ce809d62de80806577) move key functions, remove image from formation required args (@gabrtv)
- [e4da21b](https://github.com/deis/deis/commit/e4da21b26f814c0a5716103f9f1f381fe8ce716d) Revert azure to required, apparently travis CI thinks so. (@mboersma)
- [3e7ab6e](https://github.com/deis/deis/commit/3e7ab6efec8de8389f10f94b59b2c54656802f85) Fixed typo in travis after_success hook. (@mboersma)
- [5d78fb8](https://github.com/deis/deis/commit/5d78fb89838d019c1fa2629d6d328cd55f0c37ef) move default cloud init into FlavorManager (@gabrtv)
- [e421b30](https://github.com/deis/deis/commit/e421b300f1c0f75b0b6ebda7ad6fd487e47752bb) move initial_attributes into layer (@gabrtv)
- [d6ee9e9](https://github.com/deis/deis/commit/d6ee9e92ed7a19332f47fa8aae6a70da304f830e) number nodes by layer (@gabrtv)
- [bc03426](https://github.com/deis/deis/commit/bc0342674ed71b5c507728b41b7eeb080c7ac724) Filter out travis' virtualenv from coverage report. (@mboersma)
- [9bb3638](https://github.com/deis/deis/commit/9bb3638d949ce1623991c98687e119f48d8de397) add run_list and initial_attributes conditionally (@gabrtv)
- [f78a533](https://github.com/deis/deis/commit/f78a53335c394dbefc17454c0e433f637df35df3) replace formation scaling with layer & container scaling (@gabrtv)
- [b11c9f0](https://github.com/deis/deis/commit/b11c9f0fcb0a7c2f037d604afe0c2763e47acbf2) move cloud infrastructure fields from formation to layer (@gabrtv)
- [483b4d2](https://github.com/deis/deis/commit/483b4d250bb240bf2e5af7f087e4491594bc867f) include formation name in layer infrastructure (@gabrtv)
- [4a9ed1a](https://github.com/deis/deis/commit/4a9ed1af18fcd8be18d29e72a02c92dc1906741b) destroy layer infrastructure on DELETE (@gabrtv)
- [a09d60e](https://github.com/deis/deis/commit/a09d60e7d8682c4f0d3a35843d67daca3ae474dd) switch to get_object_or_404 on layer views (@gabrtv)
- [959c79c](https://github.com/deis/deis/commit/959c79c55ee27079407dd92341b40276441863d9) only apply jobs if defined (@gabrtv)
- [2cfa340](https://github.com/deis/deis/commit/2cfa3404e6005067c0f51e646618703498a4b41d) new batching logic for formation/layer destroy (@gabrtv)
- [5aed1ee](https://github.com/deis/deis/commit/5aed1ee57e3c020582c5f01b89fe76b36056a0d5) add layers:destroy, change layers:create args (@gabrtv)
- [c5638fd](https://github.com/deis/deis/commit/c5638fdb6865f8896a261c0f578a93a754d1da0b) change register url to match other auth urls (@gabrtv)
- [e115621](https://github.com/deis/deis/commit/e115621541e62433cba07a628f389553d93e71e6) remove commented line (@gabrtv)
- [3542969](https://github.com/deis/deis/commit/3542969f3fbd641349ffb2e50b153be94c1e0299) fix sg_name bug (@gabrtv)
- [fcbf263](https://github.com/deis/deis/commit/fcbf263ae0d3849e854ba57bc9a0ec4bf823c3e6) add node deletion (@gabrtv)
- [8908fd5](https://github.com/deis/deis/commit/8908fd5a72554a18172e6aba7386fba1b31bc6ad) Moved docopt to requirements.txt and updated client/Makefile. (@mboersma)
- [56f7ad1](https://github.com/deis/deis/commit/56f7ad1628c5ff788c7c67258b94bf29db69ab1d) create a list of csv run_list field (@gabrtv)
- [2a14dc7](https://github.com/deis/deis/commit/2a14dc7f80942226176a4b56b6fa1190f474b56d) Added requirements to setup.py for client. (@mboersma)
- [01df37f](https://github.com/deis/deis/commit/01df37fe52afd1d76e23c085a2384378190bd9d0) remove sleeps from node tests, test for databag on layer:scale operations (@gabrtv)
- [8f582f5](https://github.com/deis/deis/commit/8f582f5acf20a3af8233263f569d6f2659b0e74e) add chef version to layer (@gabrtv)
- [8ae3d22](https://github.com/deis/deis/commit/8ae3d22e6ba77ce71af3502f64518fefa8e57600) add Layer.level for future batching of node converges (@gabrtv)
- [5c285d0](https://github.com/deis/deis/commit/5c285d0e83a27cba267400978ccd6671471addf9) remove default build URL, allow null builds by default (@gabrtv)
- [f60c9e4](https://github.com/deis/deis/commit/f60c9e4af71eecb4930b90d9bfb54665a919dc2e) limit to just gitosis recipe on formation destroy (to assist with debugging) (@gabrtv)
- [bbfad86](https://github.com/deis/deis/commit/bbfad86fbbe9ca179b7b78dbe2bb167d180016a3) add formation info and node destroy (@gabrtv)
- [d83c2c0](https://github.com/deis/deis/commit/d83c2c040f67aac5fb663ada645f72a31426ed6f) change print on layer:create (@gabrtv)
- [663dc88](https://github.com/deis/deis/commit/663dc888cb627d824ab5dfa485045cc6f6623d96) fix bug with change flag in scale_containers (@gabrtv)
- [5b177f9](https://github.com/deis/deis/commit/5b177f900912c0375345de5f11419e9492dedc24) fix bug in layer scaling (@gabrtv)
- [c672ec4](https://github.com/deis/deis/commit/c672ec4d4984117729ee8eb827dd17e32e14f405) purge backends/proxies references from client (@gabrtv)
- [0187064](https://github.com/deis/deis/commit/01870643b780f0b2250f6149d08e918cfb0d9281) merging layout changes with current master re #34 (@gabrtv)
- [7b86a7c](https://github.com/deis/deis/commit/7b86a7cfe291006622ea33a90ed4466b0fb194eb) reset south migrations and fix some uniqueness contraints #34 (@gabrtv)
- [ea8e7ba](https://github.com/deis/deis/commit/ea8e7bac112d40f9b6ea1849125592a26dee8b47) fix client merge conflicts (@gabrtv)
- [10d1429](https://github.com/deis/deis/commit/10d1429605b8d9ea961e338dd01de645343ac58b) re-enable south, note this requires a db wipe (@gabrtv)
- [9e07a52](https://github.com/deis/deis/commit/9e07a52a5e25d059660976bc502dd0c86d848346) Build.push classmethod for handling git-push through gitosis #35 (@gabrtv)
- [519969b](https://github.com/deis/deis/commit/519969b81107e174805f561b6c89e4676cd4b2ab) rename to push-hook #35 (@gabrtv)
- [ce3430a](https://github.com/deis/deis/commit/ce3430a2a4943beeaa6a9b5d676e1937bd36df26) remove build version so push-hook doesn't have to increment it re #35 (@gabrtv)
- [4b75903](https://github.com/deis/deis/commit/4b75903384afe756072067f96d4470763166b627) catch layer does not exist error (@gabrtv)
- [b636f0d](https://github.com/deis/deis/commit/b636f0d9e0a5301f06d2434556fd24955c8c7ad1) make layer:destroy message consistent (@gabrtv)
- [e0b3724](https://github.com/deis/deis/commit/e0b3724d3ea0c9af45be4a1aed0a3c3bbac6bc07) deprecate Access and Event models, we'll reintroduce what we need later #37 (@gabrtv)
- [8efb373](https://github.com/deis/deis/commit/8efb3734606cf263fc5f125d13ad70b980055e04) purge cruft from django settings #37 (@gabrtv)
- [5b5bb71](https://github.com/deis/deis/commit/5b5bb71e1d257d70e702912ffee53999917472af) purge admin module, we'll re-enable it when we get to #38 (@gabrtv)
- [9890f0e](https://github.com/deis/deis/commit/9890f0eda97ef55719ea6e8e34895aa8a240a6f4) db migration for removing access and event models #37 (@gabrtv)
- [3370ecf](https://github.com/deis/deis/commit/3370ecf52746f693af429e6fceb925d6f4ea1f84) remove admin urls #37 (@gabrtv)
- [4f94e80](https://github.com/deis/deis/commit/4f94e80789766986e2152384ea45306976cbc86e) only save updates to formation.layers and formation.containers after successful scale operations #34 (@gabrtv)
- [0db163e](https://github.com/deis/deis/commit/0db163e6c39f0b6bcdfaac5fd0141263edbe0fc9) fix config api endpoints #16 (@gabrtv)
- [41718ff](https://github.com/deis/deis/commit/41718ffbea367b3735ef1048b5d248ece1545239) moar readme (@gabrtv)
- [8d9ec4e](https://github.com/deis/deis/commit/8d9ec4e7aed4b177ec0729b185046d5d0189b502) another round of readme updates (@gabrtv)
- [8db1ff1](https://github.com/deis/deis/commit/8db1ff1c66849bc7ab5cc77c61d2fadf674265d6) add deis-graphic (@gabrtv)
- [e6e7ff4](https://github.com/deis/deis/commit/e6e7ff4c8f496b252264647387968953cef75fb0) Updated sphinx docs layout. (@mboersma)
- [e938d57](https://github.com/deis/deis/commit/e938d57305a263a2a6c3dcad111ed89494bcf6b7) Moved "docs" target to default. (@mboersma)
- [e1c1d19](https://github.com/deis/deis/commit/e1c1d1922800a0cbe0e884d3b3b30db6d7ec1000) Added some "terms" pages to define basic Deis concepts. (@mboersma)
- [b4d32dc](https://github.com/deis/deis/commit/b4d32dc61023b34e93fb1aa6a474b83ed06bfb36) Added sphinx :ref: tags to docs. (@mboersma)
- [45d327a](https://github.com/deis/deis/commit/45d327a1a34a058908aa5bba475be9a9f047e83b) Added a sphinx theme. (@mboersma)
- [b5d515d](https://github.com/deis/deis/commit/b5d515d76431de906f86089956594a3e5a03abfc) add default run_lists to layers:create (@gabrtv)
- [4dabce9](https://github.com/deis/deis/commit/4dabce99e68116fee76e8515e47bfaa112267827) add data bag creation to provisioning (@gabrtv)
- [17409a3](https://github.com/deis/deis/commit/17409a3f958ab4821d538efae8ee4936478d735a) add provider discovery on registration #16 (@gabrtv)
- [a5d07cd](https://github.com/deis/deis/commit/a5d07cd92c400d7ac1375dec88ab70e6c0004749) create data bags and data bag items on provision (@gabrtv)
- [93c10b7](https://github.com/deis/deis/commit/93c10b75f21f1aeed2027cf2eb1754af338ec15b) update berksfile lock (@gabrtv)
- [c014426](https://github.com/deis/deis/commit/c014426bf5a2c77b4fef5a1eb3d7bb7826e03fe2) move docopt usage into docstrings and cleanup CLI dispatch re #11 #16 (@gabrtv)
- [2d5823e](https://github.com/deis/deis/commit/2d5823e25b701e001b653cf439ba1214231f542e) fix ugly git remote not found stacktrace, misc cleanup #16 (@gabrtv)
- [cb9906d](https://github.com/deis/deis/commit/cb9906da70e7ccc018e8f1a3adf713b3b52d56a3) cleanup create, scale, destroy workflow and output #16 (@gabrtv)
- [c57ae23](https://github.com/deis/deis/commit/c57ae236fa028b6fd1fd9f507010fbdd14b124d0) Refs #41 -- add client to INSTALLED_APPS, "make flake8" code cleanup. (@mboersma)
- [b04f295](https://github.com/deis/deis/commit/b04f2955b90c59fbce4d6aabccee829e34eefbb8) Added sphinx to travis ci configuration. (@mboersma)
- [02f6fad](https://github.com/deis/deis/commit/02f6fad9b2032c53311da7319e81b8fe920c5938) Fixed #41 -- repackaged client as a single-file install. (@mboersma)
- [b1f8911](https://github.com/deis/deis/commit/b1f8911fb1ea33c841ba89cb193ff920af204f5c) create one initial web container on Build.push, only if there exists a runtime layer and web containers are < 1 fixes #42 (@gabrtv)
- [914bb99](https://github.com/deis/deis/commit/914bb998f0dcc0444122d1ff0c9e079d990cf855) first pass at heroku-style container listing #16 (@gabrtv)
- [c5de935](https://github.com/deis/deis/commit/c5de935125e5dea195f6b3053dd3a6786659a2ae) resolve merge conflict on client (@gabrtv)
- [3f373db](https://github.com/deis/deis/commit/3f373db43ff4ca45a93f193c135e6a435f63a6ae) change default image to deis/buildstep #43 (@gabrtv)
- [0abdacf](https://github.com/deis/deis/commit/0abdacfb50bc67c2318971c03f99898a91c68321) order containers oldest first for CLI output #16 (@gabrtv)
- [f8a2ac1](https://github.com/deis/deis/commit/f8a2ac1fc59d03584ff1568c3bde5cffb65b9263) we need to check for > 0 runtime _nodes_, not just a runtime layer #42 (@gabrtv)
- [c7fdf8c](https://github.com/deis/deis/commit/c7fdf8c2d989ddde3dfd0df144a2df38b0a1d25d) add status fields on node/container with TODOs for adding celery beat health checks (@gabrtv)
- [de7f3fd](https://github.com/deis/deis/commit/de7f3fd66a3bd807ef172e9a7ab4040d42e89eaf) workaround for ec2 race condition (@gabrtv)
- [c3e907c](https://github.com/deis/deis/commit/c3e907cfcc6686a7713a98daba59d6d460f961a7) change managed to created, since we're not actively managing the sg and we don't want to scare admins away from locking it down (@gabrtv)
- [af4b082](https://github.com/deis/deis/commit/af4b082a58beeda8c9b13e69e9d9af814a4d5770) add subcommand help dispatch, with placeholder help for now #16 (@gabrtv)
- [eb9ad86](https://github.com/deis/deis/commit/eb9ad8658eab6ebb92497ad6cb180a5c0ba7755a) ignore node does not exist errors in the event of unclean destroy (@gabrtv)
- [45c7d00](https://github.com/deis/deis/commit/45c7d004e4744d6ea27480c9e577f735cabcd6b0) change default instance size to m1.medium (@gabrtv)
- [ed9e3a9](https://github.com/deis/deis/commit/ed9e3a9fbd5b75b4bcfaeb8ad4212674dc9c0b92) add script and instructions building Deis-optimized AMI from scratch #21 (@gabrtv)
- [ea234b5](https://github.com/deis/deis/commit/ea234b59b2de8a66093ca9cd8eb830208e6e6d29) add script to distribute AMIs across regions #21 (@gabrtv)
- [b5dda35](https://github.com/deis/deis/commit/b5dda35fdaff371f52642ce926a511419b05c642) deprecate old script (@gabrtv)
- [a61da77](https://github.com/deis/deis/commit/a61da77873dc4cce16598ca14a21513f75c2dd34) switch ec2 to new deis-optimized AMIs #21 (@gabrtv)
- [8990667](https://github.com/deis/deis/commit/8990667d55f696ee84b59181c73b6c7137954a6e) add deis-optimized amis to provision-ec2-controller script #21 (@gabrtv)
- [2079616](https://github.com/deis/deis/commit/2079616ae4bfe0bb9936155292748e56c5a17d34) minor cleanup on provision controller script (@gabrtv)
- [5d2f6a7](https://github.com/deis/deis/commit/5d2f6a7031a421ad6eff55cb1e8a97939db90d6c) check for git root before creating formation, provide better workflow guidance #16 (@gabrtv)
- [129445b](https://github.com/deis/deis/commit/129445b560e3cb099e312bcf45f573c480a4900c) cleanup formation/layer/node destroy batching (@gabrtv)
- [7d52dec](https://github.com/deis/deis/commit/7d52dec623504a382145d060c93621fcde168b69) check for no creds on layer:scale, with tests (@gabrtv)
- [08ae187](https://github.com/deis/deis/commit/08ae187fcffbcd96cc94a19c8962a3d273fb3911) only delete records from the view, fix chef_id issue (@gabrtv)
- [67c6cf3](https://github.com/deis/deis/commit/67c6cf3ce8fef854e5095eef43e21ec95eb8f4df) rework subtask batching again (@gabrtv)
- [ccd60e2](https://github.com/deis/deis/commit/ccd60e20780d32aeb5ee038598141e3adff5154c) remove unnecessary celery grouping (@gabrtv)
- [5ce8a38](https://github.com/deis/deis/commit/5ce8a38a6b2f52cc4fd719403f15c0b7c656c932) only terminate node if provider_id exists (@gabrtv)
- [6704c28](https://github.com/deis/deis/commit/6704c2805c3915b3134a369ca0ac79196b0dc944) add formation_id to args (@gabrtv)
- [31fe7f6](https://github.com/deis/deis/commit/31fe7f6ed044fc2cbce6f07739f520d85aa30449) change task invocation style (@gabrtv)
- [f79aac9](https://github.com/deis/deis/commit/f79aac9d5d4a2c5edc29ec018fa050570f0ddfda) Code cleanup via flake8. (@mboersma)
- [48fc8d1](https://github.com/deis/deis/commit/48fc8d17686f72a90746f92d4baf17c16b5a9c12) add docstrings, make `deis help <anything>` dispatch correctly #16 (@gabrtv)
- [73a1d37](https://github.com/deis/deis/commit/73a1d375191d33ae582b1b1a76f0bf6b4679462d) move parse_args() out of main, other flake8 fixes #16 (@gabrtv)
- [9a04b32](https://github.com/deis/deis/commit/9a04b32cc62b4591a0c24f5f963955d074cd0dab) resolve remaining flake8 issues (@gabrtv)
- [64deec5](https://github.com/deis/deis/commit/64deec5d3bead26d9f1fb3ec1281527717c7b14b) allow listing of builds/releases (@gabrtv)
- [6f01b25](https://github.com/deis/deis/commit/6f01b25d8061dc077ab306f41a4c797bad8a80b2) Enforce flake8 checking on travis CI. (@mboersma)
- [696b897](https://github.com/deis/deis/commit/696b8971f39b63b26d74e267ee22c30b2386032a) Updated API docs structure. (@mboersma)
- [33ddab6](https://github.com/deis/deis/commit/33ddab681c513c90aa0d3b3a25891ef878d50a2c) Remove unused "./manage.py client" command. (@mboersma)
- [ef89938](https://github.com/deis/deis/commit/ef8993858ff6923d3b31595bd79907bb4fdcc045) Added a few docstrings, refs #11. (@mboersma)
- [56d0d30](https://github.com/deis/deis/commit/56d0d30f8105d96d0331d21074d4feb65515d9e8) finish adding docstrings and inline help to cli, add support for enumerating releases/builds #16 (@gabrtv)
- [13f5128](https://github.com/deis/deis/commit/13f512881382bdfb623a5d186f237eced3c6cfae) more cli inline help edits (@gabrtv)
- [2fb15e6](https://github.com/deis/deis/commit/2fb15e6c1d332b7b75bff641bee2f25757f8b6dd) remove old readme (@gabrtv)
- [3797daf](https://github.com/deis/deis/commit/3797daf9a7778c795c9618fb62b457ed03151fa9) Refs #11 -- more docstring improvement. (@mboersma)
- [001e9d6](https://github.com/deis/deis/commit/001e9d6616cb5293c02d3c4487dd4a7197588da4) Re #11 -- more docstring progress (@mboersma)
- [578c032](https://github.com/deis/deis/commit/578c032d48b7260ddc9272e10a5620c32b3c73b6) Removed pydevd debug.py file from master branch. (@mboersma)
- [4ebd123](https://github.com/deis/deis/commit/4ebd123a1c37c539ea3e8a1909d4e5b422d8441c) standardize list and info cli output (@gabrtv)
- [abbc10a](https://github.com/deis/deis/commit/abbc10a2a0b6006353cc467cb87813293713c7df) fix flake8 line length (@gabrtv)
- [8834b94](https://github.com/deis/deis/commit/8834b946a99a40828a5e6ef01d06adab67e17e22) check for valid flavor on formations:create, add time-based done output #16 (@gabrtv)
- [72b8345](https://github.com/deis/deis/commit/72b8345ac6557a2c4707208a13a2c407478f740b) update in prep for 0.0.4 release (@gabrtv)
