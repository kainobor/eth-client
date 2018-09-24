--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: eth_balance; Type: TABLE; Schema: eth_client; Owner: postgres
--

CREATE TABLE eth_client.eth_balance (
  id integer NOT NULL,
  address character varying(42) NOT NULL,
  balance character varying(255) NOT NULL
);


ALTER TABLE eth_client.eth_balance OWNER TO postgres;

--
-- Name: TABLE eth_balance; Type: COMMENT; Schema: eth_client; Owner: postgres
--

COMMENT ON TABLE eth_client.eth_balance IS 'All known balances';


--
-- Name: eth_balance_id_seq; Type: SEQUENCE; Schema: eth_client; Owner: postgres
--

CREATE SEQUENCE eth_client.eth_balance_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE eth_client.eth_balance_id_seq OWNER TO postgres;

--
-- Name: eth_balance_id_seq; Type: SEQUENCE OWNED BY; Schema: eth_client; Owner: postgres
--

ALTER SEQUENCE eth_client.eth_balance_id_seq OWNED BY eth_client.eth_balance.id;


--
-- Name: transactions_entry; Type: TABLE; Schema: eth_client; Owner: postgres
--

CREATE TABLE eth_client.transactions_entry (
  id integer NOT NULL,
  hash character varying(66) NOT NULL,
  block_hash character varying(66),
  block_number bigint,
  from_addr character varying(42),
  to_addr character varying(42),
  confirmations integer NOT NULL,
  amount character varying(255),
  status character varying(7) DEFAULT 'pending'::character varying NOT NULL,
  showed boolean DEFAULT false NOT NULL,
  created_at timestamp without time zone
);


ALTER TABLE eth_client.transactions_entry OWNER TO postgres;

--
-- Name: transactions_entry_id_seq; Type: SEQUENCE; Schema: eth_client; Owner: postgres
--

CREATE SEQUENCE eth_client.transactions_entry_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE eth_client.transactions_entry_id_seq OWNER TO postgres;

--
-- Name: transactions_entry_id_seq; Type: SEQUENCE OWNED BY; Schema: eth_client; Owner: postgres
--

ALTER SEQUENCE eth_client.transactions_entry_id_seq OWNED BY eth_client.transactions_entry.id;


--
-- Name: transactions_withdraw; Type: TABLE; Schema: eth_client; Owner: postgres
--

CREATE TABLE eth_client.transactions_withdraw (
  id integer NOT NULL,
  hash character varying(66) NOT NULL,
  from_addr character varying(42),
  to_addr character varying(42),
  amount character varying(255),
  created_at timestamp without time zone
);


ALTER TABLE eth_client.transactions_withdraw OWNER TO postgres;

--
-- Name: transactions_withdraw_id_seq; Type: SEQUENCE; Schema: eth_client; Owner: postgres
--

CREATE SEQUENCE eth_client.transactions_withdraw_id_seq
  AS integer
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER TABLE eth_client.transactions_withdraw_id_seq OWNER TO postgres;

--
-- Name: transactions_withdraw_id_seq; Type: SEQUENCE OWNED BY; Schema: eth_client; Owner: postgres
--

ALTER SEQUENCE eth_client.transactions_withdraw_id_seq OWNED BY eth_client.transactions_withdraw.id;


--
-- Name: eth_balance id; Type: DEFAULT; Schema: eth_client; Owner: postgres
--

ALTER TABLE ONLY eth_client.eth_balance ALTER COLUMN id SET DEFAULT nextval('eth_client.eth_balance_id_seq'::regclass);


--
-- Name: transactions_entry id; Type: DEFAULT; Schema: eth_client; Owner: postgres
--

ALTER TABLE ONLY eth_client.transactions_entry ALTER COLUMN id SET DEFAULT nextval('eth_client.transactions_entry_id_seq'::regclass);


--
-- Name: transactions_withdraw id; Type: DEFAULT; Schema: eth_client; Owner: postgres
--

ALTER TABLE ONLY eth_client.transactions_withdraw ALTER COLUMN id SET DEFAULT nextval('eth_client.transactions_withdraw_id_seq'::regclass);

--
-- Name: eth_balance eth_balance_pkey; Type: CONSTRAINT; Schema: eth_client; Owner: postgres
--

ALTER TABLE ONLY eth_client.eth_balance
  ADD CONSTRAINT eth_balance_pkey PRIMARY KEY (id);


--
-- Name: transactions_entry transactions_entry_pkey; Type: CONSTRAINT; Schema: eth_client; Owner: postgres
--

ALTER TABLE ONLY eth_client.transactions_entry
  ADD CONSTRAINT transactions_entry_pkey PRIMARY KEY (id);


--
-- Name: transactions_withdraw transactions_withdraw_pkey; Type: CONSTRAINT; Schema: eth_client; Owner: postgres
--

ALTER TABLE ONLY eth_client.transactions_withdraw
  ADD CONSTRAINT transactions_withdraw_pkey PRIMARY KEY (id);


--
-- Name: balance_address_uindex; Type: INDEX; Schema: eth_client; Owner: postgres
--

CREATE UNIQUE INDEX balance_address_uindex ON eth_client.eth_balance USING btree (address);


--
-- Name: transactions_entry_hash_index; Type: INDEX; Schema: eth_client; Owner: postgres
--

CREATE INDEX transactions_entry_hash_index ON eth_client.transactions_entry USING btree (hash);


--
-- Name: transactions_withdraw_hash_index; Type: INDEX; Schema: eth_client; Owner: postgres
--

CREATE INDEX transactions_withdraw_hash_index ON eth_client.transactions_withdraw USING btree (hash);


--
-- PostgreSQL database dump complete
--
