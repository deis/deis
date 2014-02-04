--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

--
-- Data for Name: auth_user; Type: TABLE DATA; Schema: public; Owner: deis
--

COPY auth_user (id, password, last_login, is_superuser, username, first_name, last_name, email, is_staff, is_active, date_joined) FROM stdin;
1	pbkdf2_sha256$12000$18U86XLtKK5h$dV4rNDHjrkWUSPHtCdtwmSUa19OxPjOIKekvTuhs9HI=	2014-01-27 16:30:25.062197-07	t	dev			dev@dev.com	t	t	2014-01-27 15:23:49.793007-07
\.


--
-- Data for Name: api_provider; Type: TABLE DATA; Schema: public; Owner: deis
--

COPY api_provider (uuid, created, updated, owner_id, id, type, creds) FROM stdin;
74ae11e2-6080-4817-934b-826d85418ce8	2014-01-27 15:23:49.857304-07	2014-01-27 15:23:49.857328-07	1	static	static	{}
1ef7e631-79e9-40f5-9702-e2edeab2d066	2014-01-27 15:23:49.861201-07	2014-01-27 15:45:55.145015-07	1	vagrant	vagrant	{}
\.


--
-- Data for Name: api_flavor; Type: TABLE DATA; Schema: public; Owner: deis
--

COPY api_flavor (uuid, created, updated, owner_id, id, provider_id, params) FROM stdin;
4b65ecb5-3c05-4d9b-9efa-e6b8195099ce	2014-01-27 15:23:52.078972-07	2014-01-27 15:23:52.079003-07	1	vagrant-512	1ef7e631-79e9-40f5-9702-e2edeab2d066	{"memory": "512"}
ebddfd0f-9d16-42b7-9f74-ac9ba4479234	2014-01-27 15:23:52.090896-07	2014-01-27 15:23:52.09093-07	1	vagrant-1024	1ef7e631-79e9-40f5-9702-e2edeab2d066	{"memory": "1024"}
e62447cb-c08d-4e79-b3a4-ad19e52ad616	2014-01-27 15:23:52.092727-07	2014-01-27 15:23:52.092749-07	1	vagrant-2048	1ef7e631-79e9-40f5-9702-e2edeab2d066	{"memory": "2048"}
\.


--
-- Data for Name: api_formation; Type: TABLE DATA; Schema: public; Owner: deis
--

COPY api_formation (uuid, created, updated, owner_id, id, domain, nodes) FROM stdin;
8b196fb8-5774-4512-b3cc-569af8894a6c	2014-01-27 16:32:02.148594-07	2014-01-27 16:32:51.173758-07	1	dev	deis-controller.local	{}
\.


--
-- Name: auth_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: deis
--

SELECT pg_catalog.setval('auth_user_id_seq', 1, true);


--
-- PostgreSQL database dump complete
--

