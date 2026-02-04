--
-- PostgreSQL database dump
--

\restrict YuqjccSdQal1Q1f0qsAVREW1kAiV7dHPWfXFirJ5i6cvxvt0zuzr0XeXTO72nFV

-- Dumped from database version 14.20 (Homebrew)
-- Dumped by pg_dump version 14.20 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: nikunjwadhwa
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO nikunjwadhwa;

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: nikunjwadhwa
--

CREATE TABLE public.transactions (
    id uuid NOT NULL,
    amount numeric NOT NULL,
    description text NOT NULL,
    type text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.transactions OWNER TO nikunjwadhwa;

--
-- Name: users; Type: TABLE; Schema: public; Owner: nikunjwadhwa
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(255),
    password_hash character varying(255),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.users OWNER TO nikunjwadhwa;

--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: nikunjwadhwa
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: nikunjwadhwa
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: nikunjwadhwa
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: nikunjwadhwa
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: nikunjwadhwa
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- PostgreSQL database dump complete
--

\unrestrict YuqjccSdQal1Q1f0qsAVREW1kAiV7dHPWfXFirJ5i6cvxvt0zuzr0XeXTO72nFV

