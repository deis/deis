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
2	pbkdf2_sha256$12000$AMAIZeSq7IBP$6F8fYzjn7z1BBuBBDJfAA84eTW+SNknv6aqiYdc8OyY=	2014-01-20 07:58:08.834168-07	t	dev			dev@dev.com	t	t	2014-01-20 07:58:08.07811-07
\.


--
-- Data for Name: api_formation; Type: TABLE DATA; Schema: public; Owner: deis
--

COPY api_formation (uuid, created, updated, owner_id, id, domain, nodes) FROM stdin;
8edebda0-8a73-4de7-a32d-0760daa9563e	2014-01-20 07:58:10.749391-07	2014-01-20 07:58:10.74948-07	2	dev	\N	{}
\.


--
-- Name: auth_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: deis
--

SELECT pg_catalog.setval('auth_user_id_seq', 2, true);


--
-- PostgreSQL database dump complete
--

