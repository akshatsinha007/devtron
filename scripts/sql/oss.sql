--
-- PostgreSQL database dump
--

-- Dumped from database version 14.9 (Debian 14.9-1.pgdg120+1)
-- Dumped by pg_dump version 14.18 (Debian 14.18-1.pgdg120+1)

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

--
-- Name: id_seq_api_token; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_api_token
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_api_token OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: api_token; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.api_token (
    id integer DEFAULT nextval('public.id_seq_api_token'::regclass) NOT NULL,
    user_id integer NOT NULL,
    name character varying(50) NOT NULL,
    description text NOT NULL,
    expire_at_in_ms bigint,
    token text NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.api_token OWNER TO postgres;

--
-- Name: app; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app (
    id integer NOT NULL,
    app_name character varying(250) NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    team_id integer,
    app_store boolean DEFAULT false,
    app_offering_mode character varying(50) DEFAULT 'FULL'::character varying NOT NULL,
    app_type integer DEFAULT 0 NOT NULL,
    description text,
    display_name character varying(250)
);


ALTER TABLE public.app OWNER TO postgres;

--
-- Name: app_env_linkouts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_env_linkouts (
    id integer NOT NULL,
    app_id integer,
    environment_id integer,
    link text,
    description text,
    name character varying(100) NOT NULL,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer
);


ALTER TABLE public.app_env_linkouts OWNER TO postgres;

--
-- Name: app_env_linkouts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_env_linkouts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_env_linkouts_id_seq OWNER TO postgres;

--
-- Name: app_env_linkouts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_env_linkouts_id_seq OWNED BY public.app_env_linkouts.id;


--
-- Name: app_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_id_seq OWNER TO postgres;

--
-- Name: app_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_id_seq OWNED BY public.app.id;


--
-- Name: id_seq_app_label; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_app_label
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_app_label OWNER TO postgres;

--
-- Name: app_label; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_label (
    id integer DEFAULT nextval('public.id_seq_app_label'::regclass) NOT NULL,
    app_id integer NOT NULL,
    key character varying(317) NOT NULL,
    value character varying(255) NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    propagate boolean DEFAULT true NOT NULL
);


ALTER TABLE public.app_label OWNER TO postgres;

--
-- Name: app_level_metrics; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_level_metrics (
    id integer NOT NULL,
    app_id integer NOT NULL,
    app_metrics boolean NOT NULL,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    infra_metrics boolean DEFAULT true
);


ALTER TABLE public.app_level_metrics OWNER TO postgres;

--
-- Name: app_level_metrics_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_level_metrics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_level_metrics_id_seq OWNER TO postgres;

--
-- Name: app_level_metrics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_level_metrics_id_seq OWNED BY public.app_level_metrics.id;


--
-- Name: app_status; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_status (
    app_id integer NOT NULL,
    env_id integer NOT NULL,
    status character varying(50),
    updated_on timestamp with time zone NOT NULL
);


ALTER TABLE public.app_status OWNER TO postgres;

--
-- Name: app_store; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_store (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    chart_repo_id integer,
    active boolean NOT NULL,
    chart_git_location character varying(250),
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    docker_artifact_store_id character varying(250)
);


ALTER TABLE public.app_store OWNER TO postgres;

--
-- Name: app_store_application_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_store_application_version (
    id integer NOT NULL,
    version character varying(250),
    app_version character varying(250),
    created timestamp with time zone,
    deprecated boolean,
    description text,
    digest character varying(250),
    icon character varying(512),
    name character varying(256),
    home character varying(256),
    source character varying(512),
    values_yaml json NOT NULL,
    chart_yaml json NOT NULL,
    app_store_id integer,
    latest boolean DEFAULT false,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    raw_values text,
    readme text,
    created_by integer,
    updated_by integer,
    values_schema_json text,
    notes text
);


ALTER TABLE public.app_store_application_version OWNER TO postgres;

--
-- Name: app_store_application_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_store_application_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_store_application_version_id_seq OWNER TO postgres;

--
-- Name: app_store_application_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_store_application_version_id_seq OWNED BY public.app_store_application_version.id;


--
-- Name: app_store_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_store_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_store_id_seq OWNER TO postgres;

--
-- Name: app_store_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_store_id_seq OWNED BY public.app_store.id;


--
-- Name: app_store_version_values; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_store_version_values (
    id integer NOT NULL,
    name character varying(100),
    values_yaml text NOT NULL,
    app_store_application_version_id integer,
    deleted boolean DEFAULT false NOT NULL,
    created_by integer,
    updated_by integer,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    reference_type character varying(50),
    description text
);


ALTER TABLE public.app_store_version_values OWNER TO postgres;

--
-- Name: app_store_version_values_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_store_version_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_store_version_values_id_seq OWNER TO postgres;

--
-- Name: app_store_version_values_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_store_version_values_id_seq OWNED BY public.app_store_version_values.id;


--
-- Name: app_workflow; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_workflow (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    app_id integer NOT NULL,
    workflow_dag text,
    active boolean,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer
);


ALTER TABLE public.app_workflow OWNER TO postgres;

--
-- Name: app_workflow_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_workflow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_workflow_id_seq OWNER TO postgres;

--
-- Name: app_workflow_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.app_workflow_id_seq OWNED BY public.app_workflow.id;


--
-- Name: app_workflow_mapping_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.app_workflow_mapping_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.app_workflow_mapping_id_seq OWNER TO postgres;

--
-- Name: app_workflow_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.app_workflow_mapping (
    id integer DEFAULT nextval('public.app_workflow_mapping_id_seq'::regclass) NOT NULL,
    type character varying(100),
    component_id integer,
    parent_id integer,
    app_workflow_id integer NOT NULL,
    active boolean,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    parent_type character varying(100)
);


ALTER TABLE public.app_workflow_mapping OWNER TO postgres;

--
-- Name: id_artifact_promotion_approval_request; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_artifact_promotion_approval_request
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_artifact_promotion_approval_request OWNER TO postgres;

--
-- Name: artifact_promotion_approval_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.artifact_promotion_approval_request (
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    id integer DEFAULT nextval('public.id_artifact_promotion_approval_request'::regclass) NOT NULL,
    policy_id integer NOT NULL,
    policy_evaluation_audit_id integer NOT NULL,
    artifact_id integer NOT NULL,
    source_pipeline_id integer NOT NULL,
    source_type integer NOT NULL,
    destination_pipeline_id integer NOT NULL,
    status integer NOT NULL,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL
);


ALTER TABLE public.artifact_promotion_approval_request OWNER TO postgres;

--
-- Name: id_seq_attributes; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_attributes
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_attributes OWNER TO postgres;

--
-- Name: attributes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.attributes (
    id integer DEFAULT nextval('public.id_seq_attributes'::regclass) NOT NULL,
    key character varying(250) NOT NULL,
    value character varying(10000) NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.attributes OWNER TO postgres;

--
-- Name: id_seq_auto_remediation_trigger; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_auto_remediation_trigger
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_auto_remediation_trigger OWNER TO postgres;

--
-- Name: auto_remediation_trigger; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.auto_remediation_trigger (
    id integer DEFAULT nextval('public.id_seq_auto_remediation_trigger'::regclass) NOT NULL,
    type character varying(50),
    watcher_id integer,
    data text,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.auto_remediation_trigger OWNER TO postgres;

--
-- Name: id_seq_bulk_update_readme; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_bulk_update_readme
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_bulk_update_readme OWNER TO postgres;

--
-- Name: bulk_update_readme; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bulk_update_readme (
    id integer DEFAULT nextval('public.id_seq_bulk_update_readme'::regclass) NOT NULL,
    resource character varying(255) NOT NULL,
    readme text,
    script jsonb
);


ALTER TABLE public.bulk_update_readme OWNER TO postgres;

--
-- Name: casbin_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.casbin_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.casbin_id_seq OWNER TO postgres;

--
-- Name: casbin_role_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.casbin_role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.casbin_role_id_seq OWNER TO postgres;

--
-- Name: cd_workflow; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cd_workflow (
    id integer NOT NULL,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    ci_artifact_id integer NOT NULL,
    pipeline_id integer NOT NULL,
    workflow_status character varying(256)
);


ALTER TABLE public.cd_workflow OWNER TO postgres;

--
-- Name: cd_workflow_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cd_workflow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cd_workflow_id_seq OWNER TO postgres;

--
-- Name: cd_workflow_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cd_workflow_id_seq OWNED BY public.cd_workflow.id;


--
-- Name: cd_workflow_runner; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cd_workflow_runner (
    id integer NOT NULL,
    name character varying(256) NOT NULL,
    workflow_type character varying(256) NOT NULL,
    executor_type character varying(256) NOT NULL,
    status character varying(256),
    pod_status character varying(256),
    message text,
    started_on timestamp with time zone,
    finished_on timestamp with time zone,
    namespace character varying(256),
    log_file_path character varying(256),
    triggered_by integer,
    cd_workflow_id integer NOT NULL,
    blob_storage_enabled boolean DEFAULT true NOT NULL,
    pod_name text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    deployment_approval_request_id integer,
    helm_reference_chart bytea,
    ref_cd_workflow_runner_id integer,
    image_path_reservation_ids integer[],
    reference_id character varying(50),
    is_artifact_uploaded character varying(50),
    cd_artifact_location character varying(256),
    image_state character varying(50)
);


ALTER TABLE public.cd_workflow_runner OWNER TO postgres;

--
-- Name: cd_workflow_runner_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cd_workflow_runner_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cd_workflow_runner_id_seq OWNER TO postgres;

--
-- Name: cd_workflow_runner_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cd_workflow_runner_id_seq OWNED BY public.cd_workflow_runner.id;


--
-- Name: id_seq_chart_category; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_chart_category
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_chart_category OWNER TO postgres;

--
-- Name: chart_category; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_category (
    id integer DEFAULT nextval('public.id_seq_chart_category'::regclass) NOT NULL,
    name character varying(250) NOT NULL,
    description text NOT NULL,
    deleted boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.chart_category OWNER TO postgres;

--
-- Name: id_seq_chart_category_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_chart_category_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_chart_category_mapping OWNER TO postgres;

--
-- Name: chart_category_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_category_mapping (
    id integer DEFAULT nextval('public.id_seq_chart_category_mapping'::regclass) NOT NULL,
    app_store_id integer,
    chart_category_id integer,
    deleted boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.chart_category_mapping OWNER TO postgres;

--
-- Name: chart_env_config_override; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_env_config_override (
    id integer NOT NULL,
    chart_id integer,
    target_environment integer,
    env_override_yaml text NOT NULL,
    status character varying(50) NOT NULL,
    reviewed boolean NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    namespace character varying(250),
    latest boolean DEFAULT false NOT NULL,
    previous boolean DEFAULT false NOT NULL,
    is_override boolean,
    is_basic_view_locked boolean DEFAULT false NOT NULL,
    current_view_editor text DEFAULT 'UNDEFINED'::text,
    merge_strategy character varying(100)
);


ALTER TABLE public.chart_env_config_override OWNER TO postgres;

--
-- Name: chart_env_config_override_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chart_env_config_override_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chart_env_config_override_id_seq OWNER TO postgres;

--
-- Name: chart_env_config_override_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chart_env_config_override_id_seq OWNED BY public.chart_env_config_override.id;


--
-- Name: chart_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_group (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    description text,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.chart_group OWNER TO postgres;

--
-- Name: chart_group_deployment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_group_deployment (
    id integer NOT NULL,
    chart_group_id integer NOT NULL,
    chart_group_entry_id integer,
    installed_app_id integer NOT NULL,
    group_installation_id character varying(250),
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.chart_group_deployment OWNER TO postgres;

--
-- Name: chart_group_deployment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chart_group_deployment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chart_group_deployment_id_seq OWNER TO postgres;

--
-- Name: chart_group_deployment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chart_group_deployment_id_seq OWNED BY public.chart_group_deployment.id;


--
-- Name: chart_group_entry; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_group_entry (
    id integer NOT NULL,
    app_store_values_version_id integer,
    app_store_application_version_id integer,
    chart_group_id integer,
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.chart_group_entry OWNER TO postgres;

--
-- Name: chart_group_entry_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chart_group_entry_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chart_group_entry_id_seq OWNER TO postgres;

--
-- Name: chart_group_entry_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chart_group_entry_id_seq OWNED BY public.chart_group_entry.id;


--
-- Name: chart_group_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chart_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chart_group_id_seq OWNER TO postgres;

--
-- Name: chart_group_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chart_group_id_seq OWNED BY public.chart_group.id;


--
-- Name: id_seq_chart_ref; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_chart_ref
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_chart_ref OWNER TO postgres;

--
-- Name: chart_ref; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_ref (
    id integer DEFAULT nextval('public.id_seq_chart_ref'::regclass) NOT NULL,
    location character varying(250),
    version character varying(250),
    is_default boolean,
    active boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    name character varying(250) NOT NULL,
    chart_data bytea,
    chart_description text DEFAULT ''::text,
    user_uploaded boolean DEFAULT false,
    deployment_strategy_path text,
    json_path_for_strategy text,
    is_app_metrics_supported boolean DEFAULT true NOT NULL
);


ALTER TABLE public.chart_ref OWNER TO postgres;

--
-- Name: chart_ref_metadata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_ref_metadata (
    chart_name character varying(100) NOT NULL,
    chart_description text NOT NULL
);


ALTER TABLE public.chart_ref_metadata OWNER TO postgres;

--
-- Name: chart_ref_schema; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_ref_schema (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    type integer NOT NULL,
    schema text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    active boolean,
    resource_type integer NOT NULL,
    resource_value text
);


ALTER TABLE public.chart_ref_schema OWNER TO postgres;

--
-- Name: chart_ref_schema_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chart_ref_schema_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chart_ref_schema_id_seq OWNER TO postgres;

--
-- Name: chart_ref_schema_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chart_ref_schema_id_seq OWNED BY public.chart_ref_schema.id;


--
-- Name: chart_repo; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chart_repo (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    url character varying(250) NOT NULL,
    is_default boolean NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    external boolean DEFAULT false,
    user_name character varying(250),
    password character varying(250),
    ssh_key character varying(250),
    access_token character varying(250),
    auth_mode character varying(250),
    deleted boolean DEFAULT false NOT NULL,
    allow_insecure_connection boolean
);


ALTER TABLE public.chart_repo OWNER TO postgres;

--
-- Name: chart_repo_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.chart_repo_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chart_repo_id_seq OWNER TO postgres;

--
-- Name: chart_repo_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.chart_repo_id_seq OWNED BY public.chart_repo.id;


--
-- Name: charts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.charts (
    id integer NOT NULL,
    app_id integer,
    chart_repo_id integer,
    chart_name character varying(250) NOT NULL,
    chart_version character varying(250) NOT NULL,
    chart_repo character varying(250) NOT NULL,
    chart_repo_url character varying(250) NOT NULL,
    git_repo_url character varying(250),
    chart_location character varying(250),
    status character varying(50) NOT NULL,
    active boolean NOT NULL,
    reference_template character varying(250) NOT NULL,
    values_yaml text NOT NULL,
    global_override text NOT NULL,
    environment_override text,
    release_override text NOT NULL,
    user_overrides text,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    image_descriptor_template text,
    latest boolean DEFAULT false NOT NULL,
    chart_ref_id integer NOT NULL,
    pipeline_override text DEFAULT '{}'::text NOT NULL,
    previous boolean DEFAULT false NOT NULL,
    reference_chart bytea,
    is_basic_view_locked boolean DEFAULT false NOT NULL,
    current_view_editor text DEFAULT 'UNDEFINED'::text,
    is_custom_repository boolean DEFAULT false
);


ALTER TABLE public.charts OWNER TO postgres;

--
-- Name: charts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.charts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.charts_id_seq OWNER TO postgres;

--
-- Name: charts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.charts_id_seq OWNED BY public.charts.id;


--
-- Name: ci_artifact; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_artifact (
    id integer NOT NULL,
    image character varying(250),
    image_digest character varying(250),
    material_info text,
    data_source character varying(50),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    pipeline_id integer,
    ci_workflow_id integer,
    parent_ci_artifact integer,
    scan_enabled boolean DEFAULT false NOT NULL,
    scanned boolean DEFAULT false NOT NULL,
    external_ci_pipeline_id integer,
    payload_schema text,
    is_artifact_uploaded boolean DEFAULT false,
    credentials_source_type character varying(50),
    credentials_source_value character varying(50),
    component_id integer,
    target_platforms character varying(200)
);


ALTER TABLE public.ci_artifact OWNER TO postgres;

--
-- Name: ci_artifact_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_artifact_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_artifact_id_seq OWNER TO postgres;

--
-- Name: ci_artifact_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_artifact_id_seq OWNED BY public.ci_artifact.id;


--
-- Name: id_seq_ci_build_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ci_build_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ci_build_config OWNER TO postgres;

--
-- Name: ci_build_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_build_config (
    id integer DEFAULT nextval('public.id_seq_ci_build_config'::regclass) NOT NULL,
    type character varying(100),
    ci_template_id integer,
    ci_template_override_id integer,
    build_metadata text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    use_root_context boolean
);


ALTER TABLE public.ci_build_config OWNER TO postgres;

--
-- Name: id_seq_ci_env_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ci_env_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ci_env_mapping OWNER TO postgres;

--
-- Name: ci_env_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_env_mapping (
    id integer DEFAULT nextval('public.id_seq_ci_env_mapping'::regclass) NOT NULL,
    environment_id integer,
    ci_pipeline_id integer,
    deleted boolean DEFAULT false NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.ci_env_mapping OWNER TO postgres;

--
-- Name: id_seq_ci_env_mapping_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ci_env_mapping_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ci_env_mapping_history OWNER TO postgres;

--
-- Name: ci_env_mapping_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_env_mapping_history (
    id integer DEFAULT nextval('public.id_seq_ci_env_mapping_history'::regclass) NOT NULL,
    ci_pipeline_id integer,
    environment_id integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.ci_env_mapping_history OWNER TO postgres;

--
-- Name: ci_pipeline; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_pipeline (
    id integer NOT NULL,
    app_id integer,
    ci_template_id integer,
    name character varying(250),
    version character varying(250),
    active boolean NOT NULL,
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    manual boolean DEFAULT false NOT NULL,
    external boolean DEFAULT false,
    docker_args text,
    parent_ci_pipeline integer,
    scan_enabled boolean DEFAULT false NOT NULL,
    is_docker_config_overridden boolean DEFAULT false,
    ci_pipeline_type character varying(75),
    workflow_cache_config character varying(50)
);


ALTER TABLE public.ci_pipeline OWNER TO postgres;

--
-- Name: id_seq_ci_pipeline_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ci_pipeline_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ci_pipeline_history OWNER TO postgres;

--
-- Name: ci_pipeline_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_pipeline_history (
    id integer DEFAULT nextval('public.id_seq_ci_pipeline_history'::regclass) NOT NULL,
    ci_pipeline_id integer,
    ci_template_override_history text,
    ci_pipeline_material_history text,
    scan_enabled boolean,
    manual boolean,
    trigger character varying(100)
);


ALTER TABLE public.ci_pipeline_history OWNER TO postgres;

--
-- Name: ci_pipeline_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_pipeline_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_pipeline_id_seq OWNER TO postgres;

--
-- Name: ci_pipeline_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_pipeline_id_seq OWNED BY public.ci_pipeline.id;


--
-- Name: ci_pipeline_material; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_pipeline_material (
    id integer NOT NULL,
    git_material_id integer,
    ci_pipeline_id integer,
    path character varying(250),
    checkout_path character varying(250),
    type character varying(250),
    value character varying(250),
    scm_id character varying(250),
    scm_name character varying(250),
    scm_version character varying(250),
    active boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    regex character varying(50) DEFAULT ''::character varying
);


ALTER TABLE public.ci_pipeline_material OWNER TO postgres;

--
-- Name: ci_pipeline_material_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_pipeline_material_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_pipeline_material_id_seq OWNER TO postgres;

--
-- Name: ci_pipeline_material_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_pipeline_material_id_seq OWNED BY public.ci_pipeline_material.id;


--
-- Name: ci_pipeline_scripts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_pipeline_scripts (
    id integer NOT NULL,
    name character varying(256) NOT NULL,
    index integer NOT NULL,
    ci_pipeline_id integer NOT NULL,
    script text,
    stage character varying(256),
    output_location character varying(256),
    active boolean,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer
);


ALTER TABLE public.ci_pipeline_scripts OWNER TO postgres;

--
-- Name: ci_pipeline_scripts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_pipeline_scripts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_pipeline_scripts_id_seq OWNER TO postgres;

--
-- Name: ci_pipeline_scripts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_pipeline_scripts_id_seq OWNED BY public.ci_pipeline_scripts.id;


--
-- Name: ci_template; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_template (
    id integer NOT NULL,
    app_id integer,
    docker_registry_id character varying(250),
    docker_repository character varying(250),
    dockerfile_path character varying(250),
    args text,
    before_docker_build text,
    after_docker_build text,
    template_name character varying(250),
    version character varying(250),
    active boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    git_material_id integer,
    target_platform character varying(1000) DEFAULT ''::character varying NOT NULL,
    docker_build_options text,
    ci_build_config_id integer,
    build_context_git_material_id integer
);


ALTER TABLE public.ci_template OWNER TO postgres;

--
-- Name: id_seq_ci_template_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ci_template_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ci_template_history OWNER TO postgres;

--
-- Name: ci_template_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_template_history (
    id integer DEFAULT nextval('public.id_seq_ci_template_history'::regclass) NOT NULL,
    ci_template_id integer,
    app_id integer,
    docker_registry_id character varying(250),
    docker_repository character varying(250),
    dockerfile_path character varying(250),
    args text,
    before_docker_build text,
    after_docker_build text,
    template_name character varying(250),
    version character varying(250),
    target_platform character varying(1000) DEFAULT ''::character varying NOT NULL,
    docker_build_options text,
    active boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    git_material_id integer,
    ci_build_config_id integer,
    build_meta_data_type character varying(100),
    build_metadata text,
    trigger character varying(100)
);


ALTER TABLE public.ci_template_history OWNER TO postgres;

--
-- Name: ci_template_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_template_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_template_id_seq OWNER TO postgres;

--
-- Name: ci_template_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_template_id_seq OWNED BY public.ci_template.id;


--
-- Name: id_seq_ci_template_override; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ci_template_override
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ci_template_override OWNER TO postgres;

--
-- Name: ci_template_override; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_template_override (
    id integer DEFAULT nextval('public.id_seq_ci_template_override'::regclass) NOT NULL,
    ci_pipeline_id integer,
    docker_registry_id text,
    docker_repository text,
    dockerfile_path text,
    git_material_id integer,
    active boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    ci_build_config_id integer,
    build_context_git_material_id integer
);


ALTER TABLE public.ci_template_override OWNER TO postgres;

--
-- Name: ci_workflow; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ci_workflow (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    status character varying(50),
    pod_status character varying(50),
    message text,
    started_on timestamp with time zone,
    finished_on timestamp with time zone,
    namespace character varying(250),
    log_file_path character varying(250),
    git_triggers json,
    triggered_by integer NOT NULL,
    ci_pipeline_id integer NOT NULL,
    ci_artifact_location character varying(256),
    blob_storage_enabled boolean DEFAULT true NOT NULL,
    pod_name text,
    ci_build_type character varying(100),
    environment_id integer,
    ref_ci_workflow_id integer,
    parent_ci_workflow_id integer,
    image_path_reservation_id integer,
    executor_type character varying(50),
    image_path_reservation_ids integer[],
    is_artifact_uploaded character varying(50)
);


ALTER TABLE public.ci_workflow OWNER TO postgres;

--
-- Name: ci_workflow_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ci_workflow_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ci_workflow_id_seq OWNER TO postgres;

--
-- Name: ci_workflow_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ci_workflow_id_seq OWNED BY public.ci_workflow.id;


--
-- Name: cluster; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cluster (
    id integer NOT NULL,
    cluster_name character varying(250) NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    server_url character varying(250),
    config json,
    prometheus_endpoint character varying(250),
    cd_argo_setup boolean DEFAULT false,
    p_username character varying(250),
    p_password character varying(250),
    p_tls_client_cert text,
    p_tls_client_key text,
    agent_installation_stage integer DEFAULT 0,
    k8s_version character varying(250),
    error_in_connecting text,
    is_virtual_cluster boolean,
    insecure_skip_tls_verify boolean,
    proxy_url text,
    to_connect_with_ssh_tunnel boolean,
    ssh_tunnel_user character varying(100),
    ssh_tunnel_password text,
    ssh_tunnel_auth_key text,
    ssh_tunnel_server_address character varying(250),
    description text,
    remote_connection_config_id integer,
    is_prod boolean DEFAULT false
);


ALTER TABLE public.cluster OWNER TO postgres;

--
-- Name: cluster_accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cluster_accounts (
    id integer NOT NULL,
    account character varying(250) NOT NULL,
    config json NOT NULL,
    cluster_id integer NOT NULL,
    namespace character varying(250) NOT NULL,
    is_default boolean DEFAULT false,
    active boolean DEFAULT true NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.cluster_accounts OWNER TO postgres;

--
-- Name: cluster_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cluster_accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cluster_accounts_id_seq OWNER TO postgres;

--
-- Name: cluster_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cluster_accounts_id_seq OWNED BY public.cluster_accounts.id;


--
-- Name: cluster_helm_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cluster_helm_config (
    id integer NOT NULL,
    cluster_id integer NOT NULL,
    tiller_url character varying(250),
    tiller_cert character varying,
    tiller_key character varying,
    active boolean DEFAULT true NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.cluster_helm_config OWNER TO postgres;

--
-- Name: cluster_helm_config_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cluster_helm_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cluster_helm_config_id_seq OWNER TO postgres;

--
-- Name: cluster_helm_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cluster_helm_config_id_seq OWNED BY public.cluster_helm_config.id;


--
-- Name: cluster_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cluster_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cluster_id_seq OWNER TO postgres;

--
-- Name: cluster_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cluster_id_seq OWNED BY public.cluster.id;


--
-- Name: cluster_installed_apps_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cluster_installed_apps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cluster_installed_apps_id_seq OWNER TO postgres;

--
-- Name: cluster_installed_apps; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cluster_installed_apps (
    id integer DEFAULT nextval('public.cluster_installed_apps_id_seq'::regclass) NOT NULL,
    cluster_id integer,
    installed_app_id integer,
    created_by integer,
    created_on timestamp with time zone,
    updated_by integer,
    updated_on timestamp with time zone
);


ALTER TABLE public.cluster_installed_apps OWNER TO postgres;

--
-- Name: id_seq_config_map_app_level; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_config_map_app_level
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_config_map_app_level OWNER TO postgres;

--
-- Name: config_map_app_level; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.config_map_app_level (
    id integer DEFAULT nextval('public.id_seq_config_map_app_level'::regclass),
    app_id integer NOT NULL,
    config_map_data text,
    secret_data text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.config_map_app_level OWNER TO postgres;

--
-- Name: id_seq_config_map_env_level; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_config_map_env_level
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_config_map_env_level OWNER TO postgres;

--
-- Name: config_map_env_level; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.config_map_env_level (
    id integer DEFAULT nextval('public.id_seq_config_map_env_level'::regclass),
    app_id integer NOT NULL,
    environment_id integer NOT NULL,
    config_map_data text,
    secret_data text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    deleted boolean
);


ALTER TABLE public.config_map_env_level OWNER TO postgres;

--
-- Name: id_seq_config_map_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_config_map_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_config_map_history OWNER TO postgres;

--
-- Name: config_map_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.config_map_history (
    id integer DEFAULT nextval('public.id_seq_config_map_history'::regclass) NOT NULL,
    pipeline_id integer,
    app_id integer,
    data_type character varying(255),
    data text,
    deployed boolean,
    deployed_on timestamp with time zone,
    deployed_by integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.config_map_history OWNER TO postgres;

--
-- Name: id_seq_config_map_pipeline_level; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_config_map_pipeline_level
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_config_map_pipeline_level OWNER TO postgres;

--
-- Name: config_map_pipeline_level; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.config_map_pipeline_level (
    id integer DEFAULT nextval('public.id_seq_config_map_pipeline_level'::regclass),
    app_id integer NOT NULL,
    environment_id integer NOT NULL,
    pipeline_id integer NOT NULL,
    config_map_data text,
    secret_data text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.config_map_pipeline_level OWNER TO postgres;

--
-- Name: custom_tag; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.custom_tag (
    id integer NOT NULL,
    custom_tag_format text,
    tag_pattern text,
    auto_increasing_number integer DEFAULT 0,
    entity_key integer,
    entity_value text,
    active boolean DEFAULT true,
    metadata jsonb,
    enabled boolean DEFAULT false
);


ALTER TABLE public.custom_tag OWNER TO postgres;

--
-- Name: custom_tag_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.custom_tag_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.custom_tag_id_seq OWNER TO postgres;

--
-- Name: custom_tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.custom_tag_id_seq OWNED BY public.custom_tag.id;


--
-- Name: cve_policy_control_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cve_policy_control_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cve_policy_control_id_seq OWNER TO postgres;

--
-- Name: cve_policy_control; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cve_policy_control (
    id integer DEFAULT nextval('public.cve_policy_control_id_seq'::regclass) NOT NULL,
    global boolean,
    cluster_id integer,
    env_id integer,
    app_id integer,
    cve_store_id character varying(255),
    action integer,
    severity integer,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.cve_policy_control OWNER TO postgres;

--
-- Name: cve_store; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cve_store (
    name character varying(255) NOT NULL,
    severity integer,
    package character varying(255),
    version character varying(255),
    fixed_version character varying(255),
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    standard_severity integer
);


ALTER TABLE public.cve_store OWNER TO postgres;

--
-- Name: id_seq_default_auth_policy; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_default_auth_policy
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_default_auth_policy OWNER TO postgres;

--
-- Name: default_auth_policy; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.default_auth_policy (
    id integer DEFAULT nextval('public.id_seq_default_auth_policy'::regclass) NOT NULL,
    role_type character varying(250) NOT NULL,
    policy text NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    access_type character varying(50),
    entity character varying(50)
);


ALTER TABLE public.default_auth_policy OWNER TO postgres;

--
-- Name: id_seq_default_auth_role; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_default_auth_role
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_default_auth_role OWNER TO postgres;

--
-- Name: default_auth_role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.default_auth_role (
    id integer DEFAULT nextval('public.id_seq_default_auth_role'::regclass) NOT NULL,
    role_type character varying(250) NOT NULL,
    role text NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    access_type character varying(50),
    entity character varying(50)
);


ALTER TABLE public.default_auth_role OWNER TO postgres;

--
-- Name: id_seq_default_rbac_role_data; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_default_rbac_role_data
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_default_rbac_role_data OWNER TO postgres;

--
-- Name: default_rbac_role_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.default_rbac_role_data (
    id integer DEFAULT nextval('public.id_seq_default_rbac_role_data'::regclass) NOT NULL,
    role character varying(250) NOT NULL,
    default_role_data jsonb NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    enabled boolean NOT NULL
);


ALTER TABLE public.default_rbac_role_data OWNER TO postgres;

--
-- Name: id_seq_deployment_app_migration_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_deployment_app_migration_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_deployment_app_migration_history OWNER TO postgres;

--
-- Name: deployment_app_migration_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_app_migration_history (
    id integer DEFAULT nextval('public.id_seq_deployment_app_migration_history'::regclass) NOT NULL,
    app_id integer,
    env_id integer,
    is_migration_active boolean NOT NULL,
    migrate_to text,
    migrate_from text,
    current_status integer,
    error_status integer,
    error_encountered text,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.deployment_app_migration_history OWNER TO postgres;

--
-- Name: id_seq_deployment_approval_request; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_deployment_approval_request
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_deployment_approval_request OWNER TO postgres;

--
-- Name: deployment_approval_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_approval_request (
    id integer DEFAULT nextval('public.id_seq_deployment_approval_request'::regclass) NOT NULL,
    pipeline_id integer,
    ci_artifact_id integer,
    active boolean,
    artifact_deployment_triggered boolean DEFAULT false,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.deployment_approval_request OWNER TO postgres;

--
-- Name: id_seq_deployment_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_deployment_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_deployment_config OWNER TO postgres;

--
-- Name: deployment_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_config (
    id integer DEFAULT nextval('public.id_seq_deployment_config'::regclass) NOT NULL,
    app_id integer,
    environment_id integer,
    deployment_app_type character varying(100),
    config_type character varying(100),
    repo_url character varying(250),
    repo_name character varying(200),
    active boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    release_mode character varying(256) DEFAULT 'create'::character varying,
    release_config jsonb
);


ALTER TABLE public.deployment_config OWNER TO postgres;

--
-- Name: id_seq_deployment_event; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_deployment_event
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_deployment_event OWNER TO postgres;

--
-- Name: deployment_event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_event (
    id integer DEFAULT nextval('public.id_seq_deployment_event'::regclass) NOT NULL,
    app_id integer,
    env_id integer,
    pipeline_id integer,
    cd_workflow_runner_id integer,
    event_json text NOT NULL,
    metadata text NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.deployment_event OWNER TO postgres;

--
-- Name: deployment_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_group (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    status character varying(50),
    app_count integer,
    no_of_apps text,
    environment_id integer,
    ci_pipeline_id integer,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.deployment_group OWNER TO postgres;

--
-- Name: deployment_group_app; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_group_app (
    id integer NOT NULL,
    deployment_group_id integer,
    app_id integer,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.deployment_group_app OWNER TO postgres;

--
-- Name: deployment_group_app_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.deployment_group_app_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.deployment_group_app_id_seq OWNER TO postgres;

--
-- Name: deployment_group_app_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.deployment_group_app_id_seq OWNED BY public.deployment_group_app.id;


--
-- Name: deployment_group_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.deployment_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.deployment_group_id_seq OWNER TO postgres;

--
-- Name: deployment_group_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.deployment_group_id_seq OWNED BY public.deployment_group.id;


--
-- Name: deployment_status; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_status (
    id integer NOT NULL,
    app_name character varying(250) NOT NULL,
    status character varying(50) NOT NULL,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    app_id integer,
    env_id integer
);


ALTER TABLE public.deployment_status OWNER TO postgres;

--
-- Name: deployment_status_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.deployment_status_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.deployment_status_id_seq OWNER TO postgres;

--
-- Name: deployment_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.deployment_status_id_seq OWNED BY public.deployment_status.id;


--
-- Name: id_seq_deployment_template_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_deployment_template_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_deployment_template_history OWNER TO postgres;

--
-- Name: deployment_template_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_template_history (
    id integer DEFAULT nextval('public.id_seq_deployment_template_history'::regclass) NOT NULL,
    pipeline_id integer,
    app_id integer,
    target_environment integer,
    image_descriptor_template text NOT NULL,
    template text NOT NULL,
    template_name text,
    template_version text,
    is_app_metrics_enabled boolean,
    deployed boolean,
    deployed_on timestamp with time zone,
    deployed_by integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    merge_strategy character varying(100),
    template_patch_data text,
    pipeline_ids integer[]
);


ALTER TABLE public.deployment_template_history OWNER TO postgres;

--
-- Name: id_seq_devtron_resource; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource OWNER TO postgres;

--
-- Name: devtron_resource; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource (
    id integer DEFAULT nextval('public.id_seq_devtron_resource'::regclass) NOT NULL,
    kind character varying(250) NOT NULL,
    display_name character varying(250) NOT NULL,
    icon text,
    parent_kind_id integer,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    description text,
    is_exposed boolean DEFAULT true NOT NULL
);


ALTER TABLE public.devtron_resource OWNER TO postgres;

--
-- Name: id_seq_devtron_resource_object; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource_object
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource_object OWNER TO postgres;

--
-- Name: devtron_resource_object; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_object (
    id integer DEFAULT nextval('public.id_seq_devtron_resource_object'::regclass) NOT NULL,
    old_object_id integer,
    name character varying(250),
    devtron_resource_id integer,
    devtron_resource_schema_id integer,
    object_data jsonb,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    identifier text
);


ALTER TABLE public.devtron_resource_object OWNER TO postgres;

--
-- Name: id_seq_devtron_resource_object_audit; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource_object_audit
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource_object_audit OWNER TO postgres;

--
-- Name: devtron_resource_object_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_object_audit (
    id integer DEFAULT nextval('public.id_seq_devtron_resource_object_audit'::regclass) NOT NULL,
    devtron_resource_object_id integer,
    object_data json,
    audit_operation character varying(10) NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    audit_operation_path text[]
);


ALTER TABLE public.devtron_resource_object_audit OWNER TO postgres;

--
-- Name: id_seq_devtron_resource_object_dep_relations; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource_object_dep_relations
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource_object_dep_relations OWNER TO postgres;

--
-- Name: devtron_resource_object_dep_relations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_object_dep_relations (
    id integer DEFAULT nextval('public.id_seq_devtron_resource_object_dep_relations'::regclass) NOT NULL,
    component_object_id integer,
    component_dt_res_schema_id integer,
    dependency_object_id integer,
    dependency_dt_res_schema_id integer,
    type_of_dependency character varying(50),
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    dependency_object_identifier character varying(150),
    component_object_identifier character varying(150)
);


ALTER TABLE public.devtron_resource_object_dep_relations OWNER TO postgres;

--
-- Name: id_seq_devtron_resource_schema; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource_schema
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource_schema OWNER TO postgres;

--
-- Name: devtron_resource_schema; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_schema (
    id integer DEFAULT nextval('public.id_seq_devtron_resource_schema'::regclass) NOT NULL,
    devtron_resource_id integer,
    version character varying(10) NOT NULL,
    schema jsonb,
    latest boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    sample_schema json
);


ALTER TABLE public.devtron_resource_schema OWNER TO postgres;

--
-- Name: id_seq_devtron_resource_schema_audit; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource_schema_audit
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource_schema_audit OWNER TO postgres;

--
-- Name: devtron_resource_schema_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_schema_audit (
    id integer DEFAULT nextval('public.id_seq_devtron_resource_schema_audit'::regclass) NOT NULL,
    devtron_resource_schema_id integer,
    schema json,
    audit_operation character varying(10) NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.devtron_resource_schema_audit OWNER TO postgres;

--
-- Name: id_seq_devtron_resource_searchable_key; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_devtron_resource_searchable_key
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_devtron_resource_searchable_key OWNER TO postgres;

--
-- Name: devtron_resource_searchable_key; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_searchable_key (
    id integer DEFAULT nextval('public.id_seq_devtron_resource_searchable_key'::regclass) NOT NULL,
    name character varying(100),
    is_removed boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.devtron_resource_searchable_key OWNER TO postgres;

--
-- Name: id_devtron_resource_task_run; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_devtron_resource_task_run
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_devtron_resource_task_run OWNER TO postgres;

--
-- Name: devtron_resource_task_run; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.devtron_resource_task_run (
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    id integer DEFAULT nextval('public.id_devtron_resource_task_run'::regclass) NOT NULL,
    task_json jsonb NOT NULL,
    run_source_identifier character varying(500) NOT NULL,
    run_source_dependency_identifier character varying(500),
    run_target_identifier character varying(500) NOT NULL,
    task_type character varying(100) NOT NULL,
    task_type_identifier integer
);


ALTER TABLE public.devtron_resource_task_run OWNER TO postgres;

--
-- Name: docker_artifact_store; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.docker_artifact_store (
    id character varying(250) NOT NULL,
    plugin_id character varying(250) NOT NULL,
    registry_url character varying(250),
    registry_type character varying(250) NOT NULL,
    aws_accesskey_id character varying(250),
    aws_secret_accesskey character varying(250),
    aws_region character varying(250),
    username character varying(250),
    password character varying(5000),
    is_default boolean NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    connection character varying(250),
    cert text,
    is_oci_compliant_registry boolean,
    remote_connection_config_id integer,
    credentials_type character varying(124)
);


ALTER TABLE public.docker_artifact_store OWNER TO postgres;

--
-- Name: id_seq_docker_registry_ips_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_docker_registry_ips_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_docker_registry_ips_config OWNER TO postgres;

--
-- Name: docker_registry_ips_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.docker_registry_ips_config (
    id integer DEFAULT nextval('public.id_seq_docker_registry_ips_config'::regclass) NOT NULL,
    docker_artifact_store_id character varying(250) NOT NULL,
    credential_type character varying(50) NOT NULL,
    credential_value text,
    applied_cluster_ids_csv character varying(256),
    ignored_cluster_ids_csv character varying(256),
    active boolean DEFAULT true
);


ALTER TABLE public.docker_registry_ips_config OWNER TO postgres;

--
-- Name: id_seq_draft; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_draft
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_draft OWNER TO postgres;

--
-- Name: draft; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.draft (
    id integer DEFAULT nextval('public.id_seq_draft'::regclass) NOT NULL,
    app_id integer NOT NULL,
    env_id integer NOT NULL,
    resource integer NOT NULL,
    resource_name character varying(300) NOT NULL,
    draft_state integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.draft OWNER TO postgres;

--
-- Name: id_seq_draft_version; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_draft_version
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_draft_version OWNER TO postgres;

--
-- Name: draft_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.draft_version (
    id integer DEFAULT nextval('public.id_seq_draft_version'::regclass) NOT NULL,
    draft_id integer NOT NULL,
    data text NOT NULL,
    action integer NOT NULL,
    user_id integer NOT NULL,
    created_on timestamp with time zone
);


ALTER TABLE public.draft_version OWNER TO postgres;

--
-- Name: id_seq_draft_version_comment; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_draft_version_comment
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_draft_version_comment OWNER TO postgres;

--
-- Name: draft_version_comment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.draft_version_comment (
    id integer DEFAULT nextval('public.id_seq_draft_version_comment'::regclass) NOT NULL,
    draft_id integer NOT NULL,
    draft_version_id integer NOT NULL,
    comment text,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.draft_version_comment OWNER TO postgres;

--
-- Name: env_level_app_metrics; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.env_level_app_metrics (
    id integer NOT NULL,
    app_id integer NOT NULL,
    env_id integer NOT NULL,
    app_metrics boolean,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    infra_metrics boolean DEFAULT true
);


ALTER TABLE public.env_level_app_metrics OWNER TO postgres;

--
-- Name: env_level_app_metrics_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.env_level_app_metrics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.env_level_app_metrics_id_seq OWNER TO postgres;

--
-- Name: env_level_app_metrics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.env_level_app_metrics_id_seq OWNED BY public.env_level_app_metrics.id;


--
-- Name: environment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.environment (
    id integer NOT NULL,
    environment_name character varying(250) NOT NULL,
    cluster_id integer NOT NULL,
    active boolean DEFAULT true NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    "default" boolean DEFAULT false NOT NULL,
    namespace character varying(250),
    grafana_datasource_id integer,
    environment_identifier character varying(250) NOT NULL,
    description character varying(40),
    is_virtual_environment boolean
);


ALTER TABLE public.environment OWNER TO postgres;

--
-- Name: environment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.environment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.environment_id_seq OWNER TO postgres;

--
-- Name: environment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.environment_id_seq OWNED BY public.environment.id;


--
-- Name: id_seq_ephemeral_container; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ephemeral_container
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ephemeral_container OWNER TO postgres;

--
-- Name: ephemeral_container; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ephemeral_container (
    id integer DEFAULT nextval('public.id_seq_ephemeral_container'::regclass) NOT NULL,
    name character varying(253) NOT NULL,
    cluster_id integer NOT NULL,
    namespace character varying(250) NOT NULL,
    pod_name character varying(250) NOT NULL,
    target_container character varying(250) NOT NULL,
    config text NOT NULL,
    is_externally_created boolean DEFAULT false NOT NULL
);


ALTER TABLE public.ephemeral_container OWNER TO postgres;

--
-- Name: id_seq_ephemeral_container_actions; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_ephemeral_container_actions
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_ephemeral_container_actions OWNER TO postgres;

--
-- Name: ephemeral_container_actions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ephemeral_container_actions (
    id integer DEFAULT nextval('public.id_seq_ephemeral_container_actions'::regclass) NOT NULL,
    ephemeral_container_id integer NOT NULL,
    action_type integer DEFAULT 0 NOT NULL,
    performed_by integer NOT NULL,
    performed_at timestamp with time zone NOT NULL
);


ALTER TABLE public.ephemeral_container_actions OWNER TO postgres;

--
-- Name: event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event (
    id integer NOT NULL,
    event_type character varying(100) NOT NULL,
    description character varying(250)
);


ALTER TABLE public.event OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.event_id_seq OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.event_id_seq OWNED BY public.event.id;


--
-- Name: events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events (
    id integer NOT NULL,
    namespace character varying(250),
    kind character varying(250),
    component character varying(250),
    host character varying(250),
    reason character varying(250),
    status character varying(250),
    name character varying(250),
    message character varying(250),
    resource_revision character varying(250),
    creation_time_stamp timestamp with time zone,
    uid character varying(250),
    pipeline_name character varying(250),
    release_version character varying(250),
    created_on timestamp with time zone NOT NULL,
    created_by character varying(250) NOT NULL
);


ALTER TABLE public.events OWNER TO postgres;

--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.events_id_seq OWNER TO postgres;

--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: external_ci_pipeline; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.external_ci_pipeline (
    id integer NOT NULL,
    ci_pipeline_id integer,
    access_token character varying(256),
    active boolean,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    app_id integer
);


ALTER TABLE public.external_ci_pipeline OWNER TO postgres;

--
-- Name: external_ci_pipeline_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.external_ci_pipeline_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.external_ci_pipeline_id_seq OWNER TO postgres;

--
-- Name: external_ci_pipeline_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.external_ci_pipeline_id_seq OWNED BY public.external_ci_pipeline.id;


--
-- Name: id_seq_external_link; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_external_link
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_external_link OWNER TO postgres;

--
-- Name: external_link; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.external_link (
    id integer DEFAULT nextval('public.id_seq_external_link'::regclass) NOT NULL,
    external_link_monitoring_tool_id integer NOT NULL,
    name character varying(255) NOT NULL,
    url text,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    is_editable boolean DEFAULT false NOT NULL,
    description text
);


ALTER TABLE public.external_link OWNER TO postgres;

--
-- Name: id_seq_external_link_identifier_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_external_link_identifier_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_external_link_identifier_mapping OWNER TO postgres;

--
-- Name: external_link_identifier_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.external_link_identifier_mapping (
    id integer DEFAULT nextval('public.id_seq_external_link_identifier_mapping'::regclass) NOT NULL,
    external_link_id integer NOT NULL,
    cluster_id integer NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    type integer DEFAULT 0 NOT NULL,
    identifier character varying(255) DEFAULT ''::character varying NOT NULL,
    env_id integer DEFAULT 0 NOT NULL,
    app_id integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.external_link_identifier_mapping OWNER TO postgres;

--
-- Name: id_seq_external_link_monitoring_tool; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_external_link_monitoring_tool
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_external_link_monitoring_tool OWNER TO postgres;

--
-- Name: external_link_monitoring_tool; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.external_link_monitoring_tool (
    id integer DEFAULT nextval('public.id_seq_external_link_monitoring_tool'::regclass) NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255),
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    category integer
);


ALTER TABLE public.external_link_monitoring_tool OWNER TO postgres;

--
-- Name: id_seq_file_reference; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_file_reference
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_file_reference OWNER TO postgres;

--
-- Name: file_reference; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.file_reference (
    id integer DEFAULT nextval('public.id_seq_file_reference'::regclass) NOT NULL,
    data bytea,
    name character varying(255) NOT NULL,
    size bigint NOT NULL,
    mime_type character varying(255) NOT NULL,
    extension character varying(50) NOT NULL,
    file_type character varying(50) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.file_reference OWNER TO postgres;

--
-- Name: id_seq_generic_note; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_generic_note
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_generic_note OWNER TO postgres;

--
-- Name: generic_note; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.generic_note (
    id integer DEFAULT nextval('public.id_seq_generic_note'::regclass) NOT NULL,
    identifier integer NOT NULL,
    description text NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer,
    identifier_type integer
);


ALTER TABLE public.generic_note OWNER TO postgres;

--
-- Name: id_seq_generic_note_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_generic_note_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_generic_note_history OWNER TO postgres;

--
-- Name: generic_note_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.generic_note_history (
    id integer DEFAULT nextval('public.id_seq_generic_note_history'::regclass) NOT NULL,
    note_id integer NOT NULL,
    description text NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.generic_note_history OWNER TO postgres;

--
-- Name: git_host_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.git_host_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.git_host_id_seq OWNER TO postgres;

--
-- Name: git_host; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_host (
    id integer DEFAULT nextval('public.git_host_id_seq'::regclass) NOT NULL,
    name character varying(250) NOT NULL,
    active boolean NOT NULL,
    webhook_url character varying(500),
    webhook_secret character varying(250),
    event_type_header character varying(250),
    secret_header character varying(250),
    secret_validator character varying(250),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer,
    display_name character varying(250)
);


ALTER TABLE public.git_host OWNER TO postgres;

--
-- Name: git_material; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_material (
    id integer NOT NULL,
    app_id integer,
    git_provider_id integer,
    active boolean NOT NULL,
    name character varying(250),
    url character varying(250),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    checkout_path character varying(250),
    fetch_submodules boolean DEFAULT false NOT NULL,
    filter_pattern json DEFAULT '[]'::json
);


ALTER TABLE public.git_material OWNER TO postgres;

--
-- Name: id_seq_git_material_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_git_material_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_git_material_history OWNER TO postgres;

--
-- Name: git_material_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_material_history (
    id integer DEFAULT nextval('public.id_seq_git_material_history'::regclass) NOT NULL,
    app_id integer,
    git_provider_id integer,
    git_material_id integer,
    active boolean NOT NULL,
    name character varying(250),
    url character varying(250),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    checkout_path character varying(250),
    fetch_submodules boolean NOT NULL,
    filter_pattern json DEFAULT '[]'::json
);


ALTER TABLE public.git_material_history OWNER TO postgres;

--
-- Name: git_material_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.git_material_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.git_material_id_seq OWNER TO postgres;

--
-- Name: git_material_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.git_material_id_seq OWNED BY public.git_material.id;


--
-- Name: git_provider; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_provider (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    url character varying(250) NOT NULL,
    user_name character varying(25),
    password character varying(250),
    ssh_private_key text,
    access_token character varying(250),
    auth_mode character varying(250),
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    git_host_id integer,
    deleted boolean DEFAULT false NOT NULL,
    tls_key text,
    tls_cert text,
    ca_cert text,
    enable_tls_verification boolean
);


ALTER TABLE public.git_provider OWNER TO postgres;

--
-- Name: git_provider_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.git_provider_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.git_provider_id_seq OWNER TO postgres;

--
-- Name: git_provider_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.git_provider_id_seq OWNED BY public.git_provider.id;


--
-- Name: id_seq_git_sensor_node; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_git_sensor_node
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_git_sensor_node OWNER TO postgres;

--
-- Name: git_sensor_node; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_sensor_node (
    id integer DEFAULT nextval('public.id_seq_git_sensor_node'::regclass) NOT NULL,
    host character varying NOT NULL,
    port integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.git_sensor_node OWNER TO postgres;

--
-- Name: id_seq_git_sensor_node_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_git_sensor_node_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_git_sensor_node_mapping OWNER TO postgres;

--
-- Name: git_sensor_node_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_sensor_node_mapping (
    id integer DEFAULT nextval('public.id_seq_git_sensor_node_mapping'::regclass) NOT NULL,
    app_id integer NOT NULL,
    node_id integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.git_sensor_node_mapping OWNER TO postgres;

--
-- Name: git_web_hook; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.git_web_hook (
    id integer NOT NULL,
    ci_material_id integer NOT NULL,
    git_material_id integer NOT NULL,
    type character varying(250),
    value character varying(250),
    active boolean,
    last_seen_hash character varying(250),
    created_on timestamp with time zone
);


ALTER TABLE public.git_web_hook OWNER TO postgres;

--
-- Name: git_web_hook_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.git_web_hook_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.git_web_hook_id_seq OWNER TO postgres;

--
-- Name: git_web_hook_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.git_web_hook_id_seq OWNED BY public.git_web_hook.id;


--
-- Name: id_seq_gitops_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_gitops_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_gitops_config OWNER TO postgres;

--
-- Name: gitops_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.gitops_config (
    id integer DEFAULT nextval('public.id_seq_gitops_config'::regclass) NOT NULL,
    provider character varying(250) NOT NULL,
    username character varying(250) NOT NULL,
    token character varying(250),
    github_org_id character varying(250),
    host character varying(250),
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    gitlab_group_id character varying(250),
    azure_project character varying(250),
    bitbucket_workspace_id text,
    bitbucket_project_key text,
    email_id text,
    allow_custom_repository boolean DEFAULT false,
    ssh_key text,
    auth_mode text,
    ssh_host character varying(250),
    tls_cert text,
    tls_key text,
    ca_cert text,
    enable_tls_verification boolean
);


ALTER TABLE public.gitops_config OWNER TO postgres;

--
-- Name: id_seq_global_authorisation_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_authorisation_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_authorisation_config OWNER TO postgres;

--
-- Name: global_authorisation_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_authorisation_config (
    id integer DEFAULT nextval('public.id_seq_global_authorisation_config'::regclass) NOT NULL,
    config_type character varying(100) NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.global_authorisation_config OWNER TO postgres;

--
-- Name: id_seq_smtp_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_smtp_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_smtp_config OWNER TO postgres;

--
-- Name: global_cm_cs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_cm_cs (
    id integer DEFAULT nextval('public.id_seq_smtp_config'::regclass) NOT NULL,
    config_type text,
    name text,
    data text,
    mount_path text,
    deleted boolean DEFAULT false NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    type text,
    secret_ingestion_for character varying(50)
);


ALTER TABLE public.global_cm_cs OWNER TO postgres;

--
-- Name: id_seq_global_policy; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_policy
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_policy OWNER TO postgres;

--
-- Name: global_policy; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_policy (
    id integer DEFAULT nextval('public.id_seq_global_policy'::regclass) NOT NULL,
    name character varying(200),
    policy_of character varying(100),
    version character varying(20),
    policy_json text,
    description text,
    enabled boolean,
    deleted boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    policy_revision text
);


ALTER TABLE public.global_policy OWNER TO postgres;

--
-- Name: id_seq_global_policy_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_policy_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_policy_history OWNER TO postgres;

--
-- Name: global_policy_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_policy_history (
    id integer DEFAULT nextval('public.id_seq_global_policy_history'::regclass) NOT NULL,
    global_policy_id integer,
    history_of_action character varying(50),
    history_data jsonb,
    policy_of character varying(100),
    policy_version character varying(20),
    policy_data text,
    description text,
    enabled boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.global_policy_history OWNER TO postgres;

--
-- Name: id_seq_global_policy_searchable_field; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_policy_searchable_field
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_policy_searchable_field OWNER TO postgres;

--
-- Name: global_policy_searchable_field; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_policy_searchable_field (
    id integer DEFAULT nextval('public.id_seq_global_policy_searchable_field'::regclass) NOT NULL,
    global_policy_id integer,
    searchable_key_id integer,
    value text,
    is_regex boolean NOT NULL,
    policy_component integer,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    field_name character varying,
    value_int integer,
    value_time_stamp timestamp with time zone
);


ALTER TABLE public.global_policy_searchable_field OWNER TO postgres;

--
-- Name: id_seq_global_strategy_metadata; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_strategy_metadata
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_strategy_metadata OWNER TO postgres;

--
-- Name: global_strategy_metadata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_strategy_metadata (
    id integer DEFAULT nextval('public.id_seq_global_strategy_metadata'::regclass) NOT NULL,
    name text,
    description text,
    deleted boolean DEFAULT false NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    key character varying(250)
);


ALTER TABLE public.global_strategy_metadata OWNER TO postgres;

--
-- Name: COLUMN global_strategy_metadata.key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.global_strategy_metadata.key IS 'strategy json key';


--
-- Name: id_seq_global_strategy_metadata_chart_ref_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_strategy_metadata_chart_ref_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_strategy_metadata_chart_ref_mapping OWNER TO postgres;

--
-- Name: global_strategy_metadata_chart_ref_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_strategy_metadata_chart_ref_mapping (
    id integer DEFAULT nextval('public.id_seq_global_strategy_metadata_chart_ref_mapping'::regclass) NOT NULL,
    global_strategy_metadata_id integer,
    chart_ref_id integer,
    active boolean DEFAULT true NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    "default" boolean DEFAULT false
);


ALTER TABLE public.global_strategy_metadata_chart_ref_mapping OWNER TO postgres;

--
-- Name: id_seq_global_tag; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_tag
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_tag OWNER TO postgres;

--
-- Name: global_tag; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.global_tag (
    id integer DEFAULT nextval('public.id_seq_global_tag'::regclass) NOT NULL,
    key character varying(317) NOT NULL,
    mandatory_project_ids_csv text,
    propagate boolean,
    description text NOT NULL,
    active boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer,
    deployment_policy character varying(50) DEFAULT 'allow'::character varying NOT NULL,
    value_constraint_id integer
);


ALTER TABLE public.global_tag OWNER TO postgres;

--
-- Name: helm_values; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.helm_values (
    app_name character varying(250) NOT NULL,
    environment character varying(250) NOT NULL,
    values_yaml text NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.helm_values OWNER TO postgres;

--
-- Name: id_registry_index_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_registry_index_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_registry_index_mapping OWNER TO postgres;

--
-- Name: id_seq_app_group; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_app_group
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_app_group OWNER TO postgres;

--
-- Name: id_seq_app_group_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_app_group_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_app_group_mapping OWNER TO postgres;

--
-- Name: id_seq_app_store_charts_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_app_store_charts_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_app_store_charts_history OWNER TO postgres;

--
-- Name: id_seq_deployment_approval_user_data; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_deployment_approval_user_data
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_deployment_approval_user_data OWNER TO postgres;

--
-- Name: id_seq_global_cm_cs; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_global_cm_cs
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_global_cm_cs OWNER TO postgres;

--
-- Name: id_seq_image_comment; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_image_comment
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_image_comment OWNER TO postgres;

--
-- Name: id_seq_image_tag; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_image_tag
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_image_tag OWNER TO postgres;

--
-- Name: id_seq_image_tagging_audit; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_image_tagging_audit
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_image_tagging_audit OWNER TO postgres;

--
-- Name: id_seq_infra_config_trigger_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_infra_config_trigger_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_infra_config_trigger_history OWNER TO postgres;

--
-- Name: id_seq_infra_profile; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_infra_profile
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_infra_profile OWNER TO postgres;

--
-- Name: id_seq_infra_profile_configuration; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_infra_profile_configuration
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_infra_profile_configuration OWNER TO postgres;

--
-- Name: id_seq_infrastructure_installation; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_infrastructure_installation
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_infrastructure_installation OWNER TO postgres;

--
-- Name: id_seq_infrastructure_installation_versions; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_infrastructure_installation_versions
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_infrastructure_installation_versions OWNER TO postgres;

--
-- Name: id_seq_installed_app_version_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_installed_app_version_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_installed_app_version_history OWNER TO postgres;

--
-- Name: id_seq_intercepted_event_execution; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_intercepted_event_execution
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_intercepted_event_execution OWNER TO postgres;

--
-- Name: id_seq_k8s_event_watcher; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_k8s_event_watcher
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_k8s_event_watcher OWNER TO postgres;

--
-- Name: id_seq_k8s_resource_history_sequence; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_k8s_resource_history_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_k8s_resource_history_sequence OWNER TO postgres;

--
-- Name: id_seq_license_attributes; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_license_attributes
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_license_attributes OWNER TO postgres;

--
-- Name: id_seq_lock_configuration; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_lock_configuration
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_lock_configuration OWNER TO postgres;

--
-- Name: id_seq_module; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_module
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_module OWNER TO postgres;

--
-- Name: id_seq_module_action_audit_log; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_module_action_audit_log
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_module_action_audit_log OWNER TO postgres;

--
-- Name: id_seq_module_resource_status; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_module_resource_status
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_module_resource_status OWNER TO postgres;

--
-- Name: id_seq_notification_rule_sequence; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_notification_rule_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_notification_rule_sequence OWNER TO postgres;

--
-- Name: id_seq_oci_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_oci_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_oci_config OWNER TO postgres;

--
-- Name: id_seq_operation_audit; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_operation_audit
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_operation_audit OWNER TO postgres;

--
-- Name: id_seq_panel; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_panel
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_panel OWNER TO postgres;

--
-- Name: id_seq_pconfig; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pconfig
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pconfig OWNER TO postgres;

--
-- Name: id_seq_pipeline_stage; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_stage
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_stage OWNER TO postgres;

--
-- Name: id_seq_pipeline_stage_step; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_stage_step
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_stage_step OWNER TO postgres;

--
-- Name: id_seq_pipeline_stage_step_condition; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_stage_step_condition
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_stage_step_condition OWNER TO postgres;

--
-- Name: id_seq_pipeline_stage_step_variable; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_stage_step_variable
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_stage_step_variable OWNER TO postgres;

--
-- Name: id_seq_pipeline_status_timeline; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_status_timeline
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_status_timeline OWNER TO postgres;

--
-- Name: id_seq_pipeline_status_timeline_resources; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_status_timeline_resources
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_status_timeline_resources OWNER TO postgres;

--
-- Name: id_seq_pipeline_status_timeline_sync_detail; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_status_timeline_sync_detail
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_status_timeline_sync_detail OWNER TO postgres;

--
-- Name: id_seq_pipeline_strategy_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pipeline_strategy_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pipeline_strategy_history OWNER TO postgres;

--
-- Name: id_seq_plugin_metadata; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_metadata
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_metadata OWNER TO postgres;

--
-- Name: id_seq_plugin_parent_metadata; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_parent_metadata
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_parent_metadata OWNER TO postgres;

--
-- Name: id_seq_plugin_pipeline_script; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_pipeline_script
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_pipeline_script OWNER TO postgres;

--
-- Name: id_seq_plugin_stage_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_stage_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_stage_mapping OWNER TO postgres;

--
-- Name: id_seq_plugin_step; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_step
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_step OWNER TO postgres;

--
-- Name: id_seq_plugin_step_condition; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_step_condition
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_step_condition OWNER TO postgres;

--
-- Name: id_seq_plugin_step_variable; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_step_variable
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_step_variable OWNER TO postgres;

--
-- Name: id_seq_plugin_tag; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_tag
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_tag OWNER TO postgres;

--
-- Name: id_seq_plugin_tag_relation; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_plugin_tag_relation
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_plugin_tag_relation OWNER TO postgres;

--
-- Name: id_seq_pre_post_cd_script_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pre_post_cd_script_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pre_post_cd_script_history OWNER TO postgres;

--
-- Name: id_seq_pre_post_ci_script_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_pre_post_ci_script_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_pre_post_ci_script_history OWNER TO postgres;

--
-- Name: id_seq_profile_platform_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_profile_platform_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_profile_platform_mapping OWNER TO postgres;

--
-- Name: id_seq_push_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_push_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_push_config OWNER TO postgres;

--
-- Name: id_seq_rbac_policy_data; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_rbac_policy_data
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_rbac_policy_data OWNER TO postgres;

--
-- Name: id_seq_rbac_policy_resource_detail; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_rbac_policy_resource_detail
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_rbac_policy_resource_detail OWNER TO postgres;

--
-- Name: id_seq_rbac_role_audit; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_rbac_role_audit
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_rbac_role_audit OWNER TO postgres;

--
-- Name: id_seq_rbac_role_data; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_rbac_role_data
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_rbac_role_data OWNER TO postgres;

--
-- Name: id_seq_rbac_role_resource_detail; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_rbac_role_resource_detail
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_rbac_role_resource_detail OWNER TO postgres;

--
-- Name: id_seq_remote_connection_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_remote_connection_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_remote_connection_config OWNER TO postgres;

--
-- Name: id_seq_resource_protection; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_resource_protection
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_resource_protection OWNER TO postgres;

--
-- Name: id_seq_resource_protection_history; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_resource_protection_history
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_resource_protection_history OWNER TO postgres;

--
-- Name: id_seq_resource_qualifier_mapping_criteria; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_resource_qualifier_mapping_criteria
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_resource_qualifier_mapping_criteria OWNER TO postgres;

--
-- Name: id_seq_scan_step_condition; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_scan_step_condition
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_scan_step_condition OWNER TO postgres;

--
-- Name: id_seq_scan_step_condition_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_scan_step_condition_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_scan_step_condition_mapping OWNER TO postgres;

--
-- Name: id_seq_scan_tool_execution_history_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_scan_tool_execution_history_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_scan_tool_execution_history_mapping OWNER TO postgres;

--
-- Name: id_seq_scan_tool_metadata; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_scan_tool_metadata
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_scan_tool_metadata OWNER TO postgres;

--
-- Name: id_seq_scan_tool_step; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_scan_tool_step
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_scan_tool_step OWNER TO postgres;

--
-- Name: id_seq_script_path_arg_port_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_script_path_arg_port_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_script_path_arg_port_mapping OWNER TO postgres;

--
-- Name: id_seq_server_action_audit_log; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_server_action_audit_log
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_server_action_audit_log OWNER TO postgres;

--
-- Name: id_seq_sso_login_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_sso_login_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_sso_login_config OWNER TO postgres;

--
-- Name: id_seq_system_network_controller_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_system_network_controller_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_system_network_controller_config OWNER TO postgres;

--
-- Name: id_seq_terminal_access_templates; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_terminal_access_templates
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_terminal_access_templates OWNER TO postgres;

--
-- Name: id_seq_timeout_window_configuration; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_timeout_window_configuration
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_timeout_window_configuration OWNER TO postgres;

--
-- Name: id_seq_timeout_window_resource_mappings; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_timeout_window_resource_mappings
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_timeout_window_resource_mappings OWNER TO postgres;

--
-- Name: id_seq_user_audit; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_user_audit
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_user_audit OWNER TO postgres;

--
-- Name: id_seq_user_deployment_request_sequence; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_user_deployment_request_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_user_deployment_request_sequence OWNER TO postgres;

--
-- Name: id_seq_user_group; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_user_group
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_user_group OWNER TO postgres;

--
-- Name: id_seq_user_group_mapping; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_user_group_mapping
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_user_group_mapping OWNER TO postgres;

--
-- Name: id_seq_user_groups; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_user_groups
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_user_groups OWNER TO postgres;

--
-- Name: id_seq_user_terminal_access_data; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_user_terminal_access_data
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_user_terminal_access_data OWNER TO postgres;

--
-- Name: id_seq_value_constraint; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_value_constraint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_value_constraint OWNER TO postgres;

--
-- Name: id_seq_webhook_config; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_webhook_config
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_webhook_config OWNER TO postgres;

--
-- Name: id_seq_workflow_execution_stage; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.id_seq_workflow_execution_stage
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.id_seq_workflow_execution_stage OWNER TO postgres;

--
-- Name: image_comments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_comments (
    id integer DEFAULT nextval('public.id_seq_image_comment'::regclass) NOT NULL,
    comment character varying(500),
    artifact_id integer,
    user_id integer
);


ALTER TABLE public.image_comments OWNER TO postgres;

--
-- Name: image_path_reservation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_path_reservation (
    id integer NOT NULL,
    custom_tag_id integer,
    image_path text,
    active boolean DEFAULT true
);


ALTER TABLE public.image_path_reservation OWNER TO postgres;

--
-- Name: image_path_reservation_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.image_path_reservation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_path_reservation_id_seq OWNER TO postgres;

--
-- Name: image_path_reservation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.image_path_reservation_id_seq OWNED BY public.image_path_reservation.id;


--
-- Name: image_scan_deploy_info_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.image_scan_deploy_info_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_scan_deploy_info_id_seq OWNER TO postgres;

--
-- Name: image_scan_deploy_info; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_scan_deploy_info (
    id integer DEFAULT nextval('public.image_scan_deploy_info_id_seq'::regclass) NOT NULL,
    image_scan_execution_history_id integer[],
    scan_object_meta_id integer,
    object_type character varying(255),
    cluster_id integer,
    env_id integer,
    created_on timestamp without time zone,
    created_by integer,
    updated_on timestamp without time zone,
    updated_by integer
);


ALTER TABLE public.image_scan_deploy_info OWNER TO postgres;

--
-- Name: image_scan_execution_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.image_scan_execution_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_scan_execution_history_id_seq OWNER TO postgres;

--
-- Name: image_scan_execution_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_scan_execution_history (
    id integer DEFAULT nextval('public.image_scan_execution_history_id_seq'::regclass) NOT NULL,
    image character varying(255),
    execution_time timestamp with time zone,
    executed_by integer,
    image_hash character varying(255),
    source_metadata_json text,
    execution_history_directory_path text,
    source_type integer,
    source_sub_type integer,
    parent_id integer,
    is_latest boolean
);


ALTER TABLE public.image_scan_execution_history OWNER TO postgres;

--
-- Name: image_scan_execution_result_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.image_scan_execution_result_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_scan_execution_result_id_seq OWNER TO postgres;

--
-- Name: image_scan_execution_result; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_scan_execution_result (
    id integer DEFAULT nextval('public.image_scan_execution_result_id_seq'::regclass) NOT NULL,
    image_scan_execution_history_id integer NOT NULL,
    cve_store_name character varying(255) NOT NULL,
    scan_tool_id integer,
    package text,
    version text,
    fixed_version text,
    class text,
    type text,
    target text
);


ALTER TABLE public.image_scan_execution_result OWNER TO postgres;

--
-- Name: image_scan_object_meta_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.image_scan_object_meta_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_scan_object_meta_id_seq OWNER TO postgres;

--
-- Name: image_scan_object_meta; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_scan_object_meta (
    id integer DEFAULT nextval('public.image_scan_object_meta_id_seq'::regclass) NOT NULL,
    name character varying(255),
    type character varying(255),
    image character varying(255),
    active boolean
);


ALTER TABLE public.image_scan_object_meta OWNER TO postgres;

--
-- Name: image_tagging_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_tagging_audit (
    id integer DEFAULT nextval('public.id_seq_image_tagging_audit'::regclass) NOT NULL,
    data text,
    data_type integer,
    artifact_id integer,
    action integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.image_tagging_audit OWNER TO postgres;

--
-- Name: infra_config_trigger_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.infra_config_trigger_history (
    id integer DEFAULT nextval('public.id_seq_infra_config_trigger_history'::regclass) NOT NULL,
    key integer NOT NULL,
    value_string text,
    platform character varying(50) NOT NULL,
    workflow_id integer NOT NULL,
    workflow_type character varying(255) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.infra_config_trigger_history OWNER TO postgres;

--
-- Name: infra_profile; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.infra_profile (
    id integer DEFAULT nextval('public.id_seq_infra_profile'::regclass) NOT NULL,
    name character varying(50) NOT NULL,
    description character varying(350),
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    buildx_driver_type character varying(50) DEFAULT 'kubernetes'::character varying NOT NULL
);


ALTER TABLE public.infra_profile OWNER TO postgres;

--
-- Name: infra_profile_configuration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.infra_profile_configuration (
    id integer DEFAULT nextval('public.id_seq_infra_profile_configuration'::regclass) NOT NULL,
    key integer NOT NULL,
    value double precision,
    profile_id integer NOT NULL,
    unit integer NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    value_string text,
    platform character varying(50),
    profile_platform_mapping_id integer
);


ALTER TABLE public.infra_profile_configuration OWNER TO postgres;

--
-- Name: infrastructure_installation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.infrastructure_installation (
    id integer DEFAULT nextval('public.id_seq_infrastructure_installation'::regclass) NOT NULL,
    installation_type character varying(255),
    installed_entity_type character varying(64),
    installed_entity_id integer,
    installation_name character varying(128),
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    active boolean
);


ALTER TABLE public.infrastructure_installation OWNER TO postgres;

--
-- Name: infrastructure_installation_versions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.infrastructure_installation_versions (
    id integer DEFAULT nextval('public.id_seq_infrastructure_installation_versions'::regclass) NOT NULL,
    infrastructure_installation_id integer,
    installation_config text,
    action integer,
    apply_status character varying(100),
    apply_status_message character varying(200),
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    active boolean
);


ALTER TABLE public.infrastructure_installation_versions OWNER TO postgres;

--
-- Name: installed_app_version_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.installed_app_version_history (
    id integer DEFAULT nextval('public.id_seq_installed_app_version_history'::regclass) NOT NULL,
    installed_app_version_id integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    values_yaml_raw text,
    status character varying(100),
    updated_on timestamp with time zone,
    updated_by integer,
    git_hash character varying(255),
    started_on timestamp with time zone,
    finished_on timestamp with time zone,
    helm_release_status_config text,
    message text
);


ALTER TABLE public.installed_app_version_history OWNER TO postgres;

--
-- Name: installed_app_versions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.installed_app_versions (
    id integer NOT NULL,
    installed_app_id integer,
    app_store_application_version_id integer,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    values_yaml_raw text,
    active boolean DEFAULT true,
    reference_value_id integer,
    reference_value_kind character varying(250)
);


ALTER TABLE public.installed_app_versions OWNER TO postgres;

--
-- Name: installed_app_versions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.installed_app_versions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.installed_app_versions_id_seq OWNER TO postgres;

--
-- Name: installed_app_versions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.installed_app_versions_id_seq OWNED BY public.installed_app_versions.id;


--
-- Name: installed_apps; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.installed_apps (
    id integer NOT NULL,
    app_id integer,
    environment_id integer,
    created_by integer,
    updated_by integer,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    active boolean DEFAULT true,
    status integer DEFAULT 0 NOT NULL,
    git_ops_repo_name character varying(255),
    deployment_app_type character varying(50),
    deployment_app_delete_request boolean DEFAULT false,
    notes text,
    is_custom_repository boolean DEFAULT false,
    git_ops_repo_url character varying(255),
    is_manifest_scan_enabled boolean
);


ALTER TABLE public.installed_apps OWNER TO postgres;

--
-- Name: installed_apps_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.installed_apps_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.installed_apps_id_seq OWNER TO postgres;

--
-- Name: installed_apps_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.installed_apps_id_seq OWNED BY public.installed_apps.id;


--
-- Name: intercepted_event_execution; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.intercepted_event_execution (
    id integer DEFAULT nextval('public.id_seq_intercepted_event_execution'::regclass) NOT NULL,
    cluster_id integer,
    namespace character varying(250) NOT NULL,
    metadata text,
    search_data text,
    execution_message text,
    action character varying(20),
    involved_objects text,
    intercepted_at timestamp with time zone,
    status character varying(32),
    trigger_id integer,
    trigger_execution_id integer,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.intercepted_event_execution OWNER TO postgres;

--
-- Name: job_event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.job_event (
    id integer NOT NULL,
    event_trigger_time character varying(100) NOT NULL,
    name character varying(150) NOT NULL,
    status character varying(150) NOT NULL,
    message character varying(250),
    created_on timestamp with time zone,
    updated_on timestamp with time zone
);


ALTER TABLE public.job_event OWNER TO postgres;

--
-- Name: job_event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.job_event_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.job_event_id_seq OWNER TO postgres;

--
-- Name: job_event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.job_event_id_seq OWNED BY public.job_event.id;


--
-- Name: k8s_event_watcher; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.k8s_event_watcher (
    id integer DEFAULT nextval('public.id_seq_k8s_event_watcher'::regclass) NOT NULL,
    name character varying(50) NOT NULL,
    description text,
    filter_expression text NOT NULL,
    gvks text,
    selected_actions character varying(15)[],
    selectors text,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.k8s_event_watcher OWNER TO postgres;

--
-- Name: kubernetes_resource_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.kubernetes_resource_history (
    id integer DEFAULT nextval('public.id_seq_k8s_resource_history_sequence'::regclass) NOT NULL,
    app_id integer,
    app_name character varying(100),
    env_id integer,
    namespace character varying(100),
    resource_name character varying(100),
    kind character varying(100),
    "group" character varying(100),
    force_delete boolean,
    action_type character varying(100),
    deployment_app_type character varying(100),
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    version character varying(64),
    action_metadata character varying(512),
    resource character varying(64)
);


ALTER TABLE public.kubernetes_resource_history OWNER TO postgres;

--
-- Name: license_attributes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.license_attributes (
    id integer DEFAULT nextval('public.id_seq_license_attributes'::regclass) NOT NULL,
    key character varying(250) NOT NULL,
    value text NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.license_attributes OWNER TO postgres;

--
-- Name: lock_configuration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.lock_configuration (
    id integer DEFAULT nextval('public.id_seq_lock_configuration'::regclass) NOT NULL,
    config text,
    active boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.lock_configuration OWNER TO postgres;

--
-- Name: manifest_push_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.manifest_push_config (
    id integer DEFAULT nextval('public.id_seq_push_config'::regclass) NOT NULL,
    app_id integer,
    env_id integer,
    credentials_config text,
    chart_name character varying(100),
    chart_base_version character varying(100),
    storage_type character varying(100),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    deleted boolean
);


ALTER TABLE public.manifest_push_config OWNER TO postgres;

--
-- Name: module; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.module (
    id integer DEFAULT nextval('public.id_seq_module'::regclass) NOT NULL,
    name character varying(255) NOT NULL,
    version character varying(255) NOT NULL,
    status character varying(255) NOT NULL,
    updated_on timestamp with time zone,
    module_type character varying(30),
    enabled boolean
);


ALTER TABLE public.module OWNER TO postgres;

--
-- Name: module_action_audit_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.module_action_audit_log (
    id integer DEFAULT nextval('public.id_seq_module_action_audit_log'::regclass) NOT NULL,
    module_name character varying(255) NOT NULL,
    version character varying(255) NOT NULL,
    action character varying(255) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL
);


ALTER TABLE public.module_action_audit_log OWNER TO postgres;

--
-- Name: module_resource_status; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.module_resource_status (
    id integer DEFAULT nextval('public.id_seq_module_resource_status'::regclass) NOT NULL,
    module_id integer NOT NULL,
    "group" character varying(50) NOT NULL,
    version character varying(50) NOT NULL,
    kind character varying(50) NOT NULL,
    name character varying(250) NOT NULL,
    health_status character varying(50),
    health_message character varying(1024),
    active boolean,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone
);


ALTER TABLE public.module_resource_status OWNER TO postgres;

--
-- Name: notification_rule; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notification_rule (
    id integer DEFAULT nextval('public.id_seq_notification_rule_sequence'::regclass) NOT NULL,
    expression character varying(1000),
    condition_type integer DEFAULT 0 NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.notification_rule OWNER TO postgres;

--
-- Name: notification_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notification_settings (
    id integer NOT NULL,
    app_id integer,
    env_id integer,
    pipeline_id integer,
    pipeline_type character varying(50) NOT NULL,
    event_type_id integer NOT NULL,
    config json NOT NULL,
    view_id integer NOT NULL,
    team_id integer,
    notification_rule_id integer,
    additional_config_json character varying,
    cluster_id integer
);


ALTER TABLE public.notification_settings OWNER TO postgres;

--
-- Name: notification_settings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notification_settings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notification_settings_id_seq OWNER TO postgres;

--
-- Name: notification_settings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notification_settings_id_seq OWNED BY public.notification_settings.id;


--
-- Name: notification_settings_view; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notification_settings_view (
    id integer NOT NULL,
    config json NOT NULL,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    internal boolean
);


ALTER TABLE public.notification_settings_view OWNER TO postgres;

--
-- Name: notification_settings_view_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notification_settings_view_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notification_settings_view_id_seq OWNER TO postgres;

--
-- Name: notification_settings_view_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notification_settings_view_id_seq OWNED BY public.notification_settings_view.id;


--
-- Name: notification_templates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notification_templates (
    id integer NOT NULL,
    channel_type character varying(100) NOT NULL,
    node_type character varying(50) NOT NULL,
    event_type_id integer NOT NULL,
    template_name character varying(250) NOT NULL,
    template_payload text NOT NULL
);


ALTER TABLE public.notification_templates OWNER TO postgres;

--
-- Name: notification_templates_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notification_templates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notification_templates_id_seq OWNER TO postgres;

--
-- Name: notification_templates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notification_templates_id_seq OWNED BY public.notification_templates.id;


--
-- Name: notifier_event_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifier_event_log (
    id integer NOT NULL,
    destination character varying(250) NOT NULL,
    source_id integer,
    pipeline_type character varying(100) NOT NULL,
    event_type_id integer NOT NULL,
    correlation_id character varying(250) NOT NULL,
    payload text,
    is_notification_sent boolean NOT NULL,
    event_time timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL
);


ALTER TABLE public.notifier_event_log OWNER TO postgres;

--
-- Name: notifier_event_log_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notifier_event_log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notifier_event_log_id_seq OWNER TO postgres;

--
-- Name: notifier_event_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifier_event_log_id_seq OWNED BY public.notifier_event_log.id;


--
-- Name: oci_registry_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.oci_registry_config (
    id integer DEFAULT nextval('public.id_seq_oci_config'::regclass) NOT NULL,
    docker_artifact_store_id character varying(250) NOT NULL,
    repository_type character varying(100),
    repository_action character varying(100),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    deleted boolean,
    repository_list text,
    is_chart_pull_active boolean,
    is_public boolean DEFAULT false
);


ALTER TABLE public.oci_registry_config OWNER TO postgres;

--
-- Name: operation_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.operation_audit (
    id integer DEFAULT nextval('public.id_seq_operation_audit'::regclass) NOT NULL,
    entity_id integer NOT NULL,
    entity_type character varying(50) NOT NULL,
    operation_type character varying(20) NOT NULL,
    entity_value_json jsonb NOT NULL,
    entity_value_schema_type character varying(20) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.operation_audit OWNER TO postgres;

--
-- Name: panel; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.panel (
    id integer DEFAULT nextval('public.id_seq_panel'::regclass) NOT NULL,
    name character varying(250) NOT NULL,
    cluster_id integer NOT NULL,
    active boolean NOT NULL,
    embed_iframe text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.panel OWNER TO postgres;

--
-- Name: pipeline; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline (
    id integer NOT NULL,
    app_id integer,
    ci_pipeline_id integer,
    trigger_type character varying(250) NOT NULL,
    environment_id integer,
    deployment_template character varying(250),
    pipeline_name character varying(250) NOT NULL,
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    pipeline_override text DEFAULT '{}'::text,
    pre_stage_config_yaml text,
    post_stage_config_yaml text,
    pre_trigger_type character varying(250),
    post_trigger_type character varying(250),
    pre_stage_config_map_secret_names text,
    post_stage_config_map_secret_names text,
    run_pre_stage_in_env boolean DEFAULT false,
    run_post_stage_in_env boolean DEFAULT false,
    deployment_app_created boolean DEFAULT false NOT NULL,
    deployment_app_type character varying(50),
    deployment_app_name text,
    deployment_app_delete_request boolean DEFAULT false,
    user_approval_config character varying(1000)
);


ALTER TABLE public.pipeline OWNER TO postgres;

--
-- Name: pipeline_config_override; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_config_override (
    id integer NOT NULL,
    request_identifier character varying(250) NOT NULL,
    env_config_override_id integer,
    pipeline_override_yaml text NOT NULL,
    merged_values_yaml text NOT NULL,
    status character varying(50) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    git_hash character varying(250),
    ci_artifact_id integer,
    pipeline_id integer,
    pipeline_release_counter integer,
    cd_workflow_id integer,
    deployment_type integer DEFAULT 0,
    commit_time timestamp with time zone,
    switch_traffic boolean
);


ALTER TABLE public.pipeline_config_override OWNER TO postgres;

--
-- Name: pipeline_config_override_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pipeline_config_override_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pipeline_config_override_id_seq OWNER TO postgres;

--
-- Name: pipeline_config_override_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pipeline_config_override_id_seq OWNED BY public.pipeline_config_override.id;


--
-- Name: pipeline_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pipeline_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pipeline_id_seq OWNER TO postgres;

--
-- Name: pipeline_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pipeline_id_seq OWNED BY public.pipeline.id;


--
-- Name: pipeline_stage; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_stage (
    id integer DEFAULT nextval('public.id_seq_pipeline_stage'::regclass) NOT NULL,
    name text,
    description text,
    type character varying(255),
    deleted boolean,
    ci_pipeline_id integer,
    cd_pipeline_id integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pipeline_stage OWNER TO postgres;

--
-- Name: pipeline_stage_step; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_stage_step (
    id integer DEFAULT nextval('public.id_seq_pipeline_stage_step'::regclass) NOT NULL,
    pipeline_stage_id integer,
    name character varying(255),
    description text,
    index integer,
    step_type character varying(255),
    script_id integer,
    ref_plugin_id integer,
    output_directory_path text[],
    dependent_on_step text,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    trigger_if_parent_stage_fail boolean
);


ALTER TABLE public.pipeline_stage_step OWNER TO postgres;

--
-- Name: pipeline_stage_step_condition; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_stage_step_condition (
    id integer DEFAULT nextval('public.id_seq_pipeline_stage_step_condition'::regclass) NOT NULL,
    pipeline_stage_step_id integer,
    condition_variable_id integer,
    condition_type character varying(255),
    conditional_operator character varying(255),
    conditional_value character varying(255),
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pipeline_stage_step_condition OWNER TO postgres;

--
-- Name: pipeline_stage_step_variable; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_stage_step_variable (
    id integer DEFAULT nextval('public.id_seq_pipeline_stage_step_variable'::regclass) NOT NULL,
    pipeline_stage_step_id integer,
    name character varying(255),
    format character varying(255),
    description text,
    is_exposed boolean,
    allow_empty_value boolean,
    default_value text,
    value text,
    variable_type character varying(255),
    index integer,
    value_type character varying(255),
    previous_step_index integer,
    variable_step_index_in_plugin integer,
    reference_variable_name text,
    reference_variable_stage text,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    value_constraint_id integer,
    file_reference_id integer,
    file_mount_dir character varying(255),
    is_runtime_arg boolean DEFAULT false NOT NULL
);


ALTER TABLE public.pipeline_stage_step_variable OWNER TO postgres;

--
-- Name: pipeline_status_timeline; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_status_timeline (
    id integer DEFAULT nextval('public.id_seq_pipeline_status_timeline'::regclass) NOT NULL,
    status character varying(255),
    status_detail text,
    status_time timestamp with time zone,
    cd_workflow_runner_id integer,
    installed_app_version_history_id integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pipeline_status_timeline OWNER TO postgres;

--
-- Name: pipeline_status_timeline_resources; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_status_timeline_resources (
    id integer DEFAULT nextval('public.id_seq_pipeline_status_timeline_resources'::regclass) NOT NULL,
    installed_app_version_history_id integer,
    cd_workflow_runner_id integer,
    resource_name character varying(1000),
    resource_kind character varying(1000),
    resource_group character varying(1000),
    resource_phase text,
    resource_status text,
    status_message text,
    timeline_stage character varying(100) DEFAULT 'KUBECTL_APPLY'::character varying,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pipeline_status_timeline_resources OWNER TO postgres;

--
-- Name: pipeline_status_timeline_sync_detail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_status_timeline_sync_detail (
    id integer DEFAULT nextval('public.id_seq_pipeline_status_timeline_sync_detail'::regclass) NOT NULL,
    installed_app_version_history_id integer,
    cd_workflow_runner_id integer,
    last_synced_at timestamp with time zone,
    sync_count integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pipeline_status_timeline_sync_detail OWNER TO postgres;

--
-- Name: pipeline_strategy; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_strategy (
    id integer NOT NULL,
    strategy character varying(250) NOT NULL,
    config text,
    created_by integer,
    updated_by integer,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    deleted boolean,
    "default" boolean NOT NULL,
    pipeline_id integer NOT NULL
);


ALTER TABLE public.pipeline_strategy OWNER TO postgres;

--
-- Name: pipeline_strategy_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pipeline_strategy_history (
    id integer DEFAULT nextval('public.id_seq_pipeline_strategy_history'::regclass) NOT NULL,
    pipeline_id integer NOT NULL,
    config text,
    strategy text NOT NULL,
    "default" boolean,
    deployed boolean,
    deployed_on timestamp with time zone,
    deployed_by integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    pipeline_trigger_type character varying(255)
);


ALTER TABLE public.pipeline_strategy_history OWNER TO postgres;

--
-- Name: pipeline_strategy_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pipeline_strategy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pipeline_strategy_id_seq OWNER TO postgres;

--
-- Name: pipeline_strategy_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pipeline_strategy_id_seq OWNED BY public.pipeline_strategy.id;


--
-- Name: plugin_metadata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_metadata (
    id integer DEFAULT nextval('public.id_seq_plugin_metadata'::regclass) NOT NULL,
    name text,
    description text,
    type character varying(255),
    icon text,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    plugin_parent_metadata_id integer,
    plugin_version text DEFAULT '1.0.0'::text NOT NULL,
    is_deprecated boolean DEFAULT false NOT NULL,
    is_latest boolean DEFAULT true NOT NULL,
    doc_link text,
    is_exposed boolean DEFAULT true NOT NULL
);


ALTER TABLE public.plugin_metadata OWNER TO postgres;

--
-- Name: plugin_parent_metadata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_parent_metadata (
    id integer DEFAULT nextval('public.id_seq_plugin_parent_metadata'::regclass) NOT NULL,
    name text NOT NULL,
    identifier text NOT NULL,
    deleted boolean NOT NULL,
    description text,
    type character varying(255),
    icon text,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    is_exposed boolean DEFAULT true NOT NULL
);


ALTER TABLE public.plugin_parent_metadata OWNER TO postgres;

--
-- Name: plugin_pipeline_script; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_pipeline_script (
    id integer DEFAULT nextval('public.id_seq_plugin_pipeline_script'::regclass) NOT NULL,
    script text,
    type character varying(255),
    store_script_at text,
    dockerfile_exists boolean,
    mount_path text,
    mount_code_to_container boolean,
    mount_code_to_container_path text,
    mount_directory_from_host boolean,
    container_image_path text,
    image_pull_secret_type character varying(255),
    image_pull_secret text,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.plugin_pipeline_script OWNER TO postgres;

--
-- Name: plugin_stage_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_stage_mapping (
    id integer DEFAULT nextval('public.id_seq_plugin_stage_mapping'::regclass) NOT NULL,
    plugin_id integer NOT NULL,
    stage_type integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.plugin_stage_mapping OWNER TO postgres;

--
-- Name: plugin_step; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_step (
    id integer DEFAULT nextval('public.id_seq_plugin_step'::regclass) NOT NULL,
    plugin_id integer,
    name character varying(255),
    description text,
    index integer,
    step_type character varying(255),
    script_id integer,
    ref_plugin_id integer,
    output_directory_path text[],
    dependent_on_step text,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.plugin_step OWNER TO postgres;

--
-- Name: plugin_step_condition; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_step_condition (
    id integer DEFAULT nextval('public.id_seq_plugin_step_condition'::regclass) NOT NULL,
    plugin_step_id integer,
    condition_variable_id integer,
    condition_type character varying(255),
    conditional_operator character varying(255),
    conditional_value character varying(255),
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.plugin_step_condition OWNER TO postgres;

--
-- Name: plugin_step_variable; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_step_variable (
    id integer DEFAULT nextval('public.id_seq_plugin_step_variable'::regclass) NOT NULL,
    plugin_step_id integer,
    name character varying(255),
    format character varying(255),
    description text,
    is_exposed boolean,
    allow_empty_value boolean,
    default_value character varying(255),
    value character varying(255),
    variable_type character varying(255),
    value_type character varying(255),
    previous_step_index integer,
    variable_step_index integer,
    variable_step_index_in_plugin integer,
    reference_variable_name text,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    file_reference_id integer,
    file_mount_dir character varying(255)
);


ALTER TABLE public.plugin_step_variable OWNER TO postgres;

--
-- Name: plugin_tag; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_tag (
    id integer DEFAULT nextval('public.id_seq_plugin_tag'::regclass) NOT NULL,
    name character varying(255),
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.plugin_tag OWNER TO postgres;

--
-- Name: plugin_tag_relation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plugin_tag_relation (
    id integer DEFAULT nextval('public.id_seq_plugin_tag_relation'::regclass) NOT NULL,
    tag_id integer,
    plugin_id integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.plugin_tag_relation OWNER TO postgres;

--
-- Name: pre_post_cd_script_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pre_post_cd_script_history (
    id integer DEFAULT nextval('public.id_seq_pre_post_cd_script_history'::regclass) NOT NULL,
    pipeline_id integer NOT NULL,
    script text,
    stage text,
    configmap_secret_names text,
    configmap_data text,
    secret_data text,
    exec_in_env boolean,
    trigger_type text,
    deployed boolean,
    deployed_on timestamp with time zone,
    deployed_by integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pre_post_cd_script_history OWNER TO postgres;

--
-- Name: pre_post_ci_script_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pre_post_ci_script_history (
    id integer DEFAULT nextval('public.id_seq_pre_post_ci_script_history'::regclass) NOT NULL,
    ci_pipeline_scripts_id integer NOT NULL,
    script text,
    stage text,
    name text,
    output_location text,
    built boolean,
    built_on timestamp with time zone,
    built_by integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.pre_post_ci_script_history OWNER TO postgres;

--
-- Name: profile_platform_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.profile_platform_mapping (
    id integer DEFAULT nextval('public.id_seq_profile_platform_mapping'::regclass) NOT NULL,
    profile_id integer NOT NULL,
    platform character varying(50) NOT NULL,
    active boolean NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    created_on timestamp with time zone NOT NULL
);


ALTER TABLE public.profile_platform_mapping OWNER TO postgres;

--
-- Name: rbac_policy_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rbac_policy_data (
    id integer DEFAULT nextval('public.id_seq_rbac_policy_data'::regclass) NOT NULL,
    entity character varying(250) NOT NULL,
    access_type character varying(250) NOT NULL,
    role character varying(250) NOT NULL,
    policy_data jsonb NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    is_preset_role boolean DEFAULT false NOT NULL,
    deleted boolean NOT NULL
);


ALTER TABLE public.rbac_policy_data OWNER TO postgres;

--
-- Name: rbac_policy_resource_detail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rbac_policy_resource_detail (
    id integer DEFAULT nextval('public.id_seq_rbac_policy_resource_detail'::regclass) NOT NULL,
    resource character varying(250) NOT NULL,
    policy_resource_value jsonb,
    allowed_actions character varying(100)[],
    resource_object jsonb,
    eligible_entity_access_types character varying(250)[],
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.rbac_policy_resource_detail OWNER TO postgres;

--
-- Name: rbac_role_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rbac_role_audit (
    id integer DEFAULT nextval('public.id_seq_rbac_role_audit'::regclass) NOT NULL,
    entity character varying(250) NOT NULL,
    access_type character varying(250),
    role character varying(250) NOT NULL,
    policy_data jsonb,
    role_data jsonb,
    audit_operation character varying(20) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.rbac_role_audit OWNER TO postgres;

--
-- Name: rbac_role_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rbac_role_data (
    id integer DEFAULT nextval('public.id_seq_rbac_role_data'::regclass) NOT NULL,
    entity character varying(250) NOT NULL,
    access_type character varying(250) NOT NULL,
    role character varying(250) NOT NULL,
    role_display_name character varying(250) NOT NULL,
    role_data jsonb NOT NULL,
    role_description text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    is_preset_role boolean DEFAULT false NOT NULL,
    deleted boolean NOT NULL
);


ALTER TABLE public.rbac_role_data OWNER TO postgres;

--
-- Name: rbac_role_resource_detail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rbac_role_resource_detail (
    id integer DEFAULT nextval('public.id_seq_rbac_role_resource_detail'::regclass) NOT NULL,
    resource character varying(250) NOT NULL,
    role_resource_key character varying(100),
    role_resource_update_key character varying(100),
    eligible_entity_access_types character varying(250)[],
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    role_resource_version character varying(250)[] DEFAULT ARRAY['base'::character varying]
);


ALTER TABLE public.rbac_role_resource_detail OWNER TO postgres;

--
-- Name: registry_index_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.registry_index_mapping (
    id integer DEFAULT nextval('public.id_registry_index_mapping'::regclass) NOT NULL,
    scan_tool_id integer,
    registry_type character varying(20),
    starting_index integer
);


ALTER TABLE public.registry_index_mapping OWNER TO postgres;

--
-- Name: release_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.release_tags (
    id integer DEFAULT nextval('public.id_seq_image_tag'::regclass) NOT NULL,
    tag_name character varying(128),
    artifact_id integer,
    deleted boolean,
    app_id integer
);


ALTER TABLE public.release_tags OWNER TO postgres;

--
-- Name: remote_connection_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.remote_connection_config (
    id integer DEFAULT nextval('public.id_seq_remote_connection_config'::regclass) NOT NULL,
    connection_method character varying(50) NOT NULL,
    proxy_url character varying(300),
    ssh_server_address character varying(300),
    ssh_username character varying(300),
    ssh_password text,
    ssh_auth_key text,
    deleted boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.remote_connection_config OWNER TO postgres;

--
-- Name: request_approval_user_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.request_approval_user_data (
    id integer DEFAULT nextval('public.id_seq_deployment_approval_user_data'::regclass) NOT NULL,
    approval_request_id integer,
    user_id integer,
    user_response integer,
    comments character varying(1000),
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    request_type integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.request_approval_user_data OWNER TO postgres;

--
-- Name: resource_filter_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_filter_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resource_filter_seq OWNER TO postgres;

--
-- Name: resource_filter; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_filter (
    id integer DEFAULT nextval('public.resource_filter_seq'::regclass) NOT NULL,
    name character varying(300) NOT NULL,
    target_object integer NOT NULL,
    condition_expression text NOT NULL,
    deleted boolean NOT NULL,
    description text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.resource_filter OWNER TO postgres;

--
-- Name: resource_filter_audit_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_filter_audit_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resource_filter_audit_seq OWNER TO postgres;

--
-- Name: resource_filter_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_filter_audit (
    id integer DEFAULT nextval('public.resource_filter_audit_seq'::regclass) NOT NULL,
    target_object integer NOT NULL,
    conditions text NOT NULL,
    filter_id integer NOT NULL,
    action integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    filter_name character varying(300)
);


ALTER TABLE public.resource_filter_audit OWNER TO postgres;

--
-- Name: resource_filter_evaluation_audit_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_filter_evaluation_audit_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resource_filter_evaluation_audit_seq OWNER TO postgres;

--
-- Name: resource_filter_evaluation_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_filter_evaluation_audit (
    id integer DEFAULT nextval('public.resource_filter_evaluation_audit_seq'::regclass) NOT NULL,
    reference_type integer NOT NULL,
    reference_id integer NOT NULL,
    filter_history_objects text NOT NULL,
    subject_type integer NOT NULL,
    subject_id integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    filter_type integer DEFAULT 1
);


ALTER TABLE public.resource_filter_evaluation_audit OWNER TO postgres;

--
-- Name: resource_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_group (
    id integer DEFAULT nextval('public.id_seq_app_group'::regclass) NOT NULL,
    name character varying(250) NOT NULL,
    description character varying(50),
    resource_id integer NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    resource_key integer DEFAULT 7 NOT NULL
);


ALTER TABLE public.resource_group OWNER TO postgres;

--
-- Name: resource_group_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_group_mapping (
    id integer DEFAULT nextval('public.id_seq_app_group_mapping'::regclass) NOT NULL,
    resource_group_id integer NOT NULL,
    resource_id integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    resource_key integer DEFAULT 6 NOT NULL
);


ALTER TABLE public.resource_group_mapping OWNER TO postgres;

--
-- Name: resource_protection; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_protection (
    id integer DEFAULT nextval('public.id_seq_resource_protection'::regclass) NOT NULL,
    app_id integer NOT NULL,
    env_id integer NOT NULL,
    resource integer NOT NULL,
    protection_state integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.resource_protection OWNER TO postgres;

--
-- Name: resource_protection_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_protection_history (
    id integer DEFAULT nextval('public.id_seq_resource_protection_history'::regclass) NOT NULL,
    app_id integer NOT NULL,
    env_id integer NOT NULL,
    resource integer NOT NULL,
    protection_state integer NOT NULL,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.resource_protection_history OWNER TO postgres;

--
-- Name: variable_scope_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variable_scope_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variable_scope_seq OWNER TO postgres;

--
-- Name: resource_qualifier_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_qualifier_mapping (
    id integer DEFAULT nextval('public.variable_scope_seq'::regclass) NOT NULL,
    resource_id integer,
    qualifier_id integer,
    identifier_key integer,
    identifier_value_int integer,
    identifier_value_string character varying(255),
    active boolean NOT NULL,
    parent_identifier integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    resource_type integer DEFAULT 0,
    global_policy_id integer
);


ALTER TABLE public.resource_qualifier_mapping OWNER TO postgres;

--
-- Name: resource_qualifier_mapping_criteria; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_qualifier_mapping_criteria (
    id integer DEFAULT nextval('public.id_seq_resource_qualifier_mapping_criteria'::regclass) NOT NULL,
    description character varying(100),
    json_data text,
    active boolean,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL,
    policy_type character varying(50) NOT NULL,
    version character varying(50) DEFAULT 'v1'::character varying NOT NULL
);


ALTER TABLE public.resource_qualifier_mapping_criteria OWNER TO postgres;

--
-- Name: resource_scan_execution_result_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_scan_execution_result_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resource_scan_execution_result_id_seq OWNER TO postgres;

--
-- Name: resource_scan_execution_result; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_scan_execution_result (
    id integer DEFAULT nextval('public.resource_scan_execution_result_id_seq'::regclass) NOT NULL,
    image_scan_execution_history_id integer NOT NULL,
    scan_data_json text,
    format integer,
    types integer[],
    scan_tool_id integer
);


ALTER TABLE public.resource_scan_execution_result OWNER TO postgres;

--
-- Name: role_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_group (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    casbin_name character varying(100),
    description text,
    created_by integer,
    updated_by integer,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    active boolean DEFAULT true NOT NULL
);


ALTER TABLE public.role_group OWNER TO postgres;

--
-- Name: role_group_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.role_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.role_group_id_seq OWNER TO postgres;

--
-- Name: role_group_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.role_group_id_seq OWNED BY public.role_group.id;


--
-- Name: role_group_role_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_group_role_mapping (
    id integer NOT NULL,
    role_group_id integer NOT NULL,
    role_id integer NOT NULL,
    created_by integer,
    updated_by integer,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone NOT NULL
);


ALTER TABLE public.role_group_role_mapping OWNER TO postgres;

--
-- Name: role_group_role_mapping_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.role_group_role_mapping_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.role_group_role_mapping_id_seq OWNER TO postgres;

--
-- Name: role_group_role_mapping_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.role_group_role_mapping_id_seq OWNED BY public.role_group_role_mapping.id;


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.roles_id_seq OWNER TO postgres;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id integer DEFAULT nextval('public.roles_id_seq'::regclass) NOT NULL,
    role character varying(250) NOT NULL,
    team character varying(100),
    environment text,
    entity_name text,
    action character varying(100),
    created_by integer,
    created_on timestamp without time zone,
    updated_by integer,
    updated_on timestamp without time zone,
    entity character varying(100),
    access_type character varying(100),
    cluster text,
    namespace text,
    "group" text,
    kind text,
    resource text,
    approver boolean,
    workflow text,
    release text,
    release_track text,
    subaction character varying(100)
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: scan_step_condition; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.scan_step_condition (
    id integer DEFAULT nextval('public.id_seq_scan_step_condition'::regclass) NOT NULL,
    condition_variable_format character varying(10),
    conditional_operator character varying(5),
    conditional_value character varying(100),
    condition_on character varying(20),
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.scan_step_condition OWNER TO postgres;

--
-- Name: scan_step_condition_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.scan_step_condition_mapping (
    id integer DEFAULT nextval('public.id_seq_scan_step_condition_mapping'::regclass) NOT NULL,
    scan_step_condition_id integer,
    scan_tool_step_id integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.scan_step_condition_mapping OWNER TO postgres;

--
-- Name: scan_tool_execution_history_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.scan_tool_execution_history_mapping (
    id integer DEFAULT nextval('public.id_seq_scan_tool_execution_history_mapping'::regclass) NOT NULL,
    image_scan_execution_history_id integer,
    scan_tool_id integer,
    execution_start_time timestamp with time zone,
    execution_finish_time timestamp with time zone,
    state integer,
    try_count integer,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    error_message character varying
);


ALTER TABLE public.scan_tool_execution_history_mapping OWNER TO postgres;

--
-- Name: scan_tool_metadata; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.scan_tool_metadata (
    id integer DEFAULT nextval('public.id_seq_scan_tool_metadata'::regclass) NOT NULL,
    name character varying(100),
    version character varying(50),
    server_base_url character varying(30),
    result_descriptor_template text,
    scan_target character varying(10),
    active boolean,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    tool_metadata text,
    is_preset boolean DEFAULT true NOT NULL,
    plugin_id integer,
    url character varying(100)
);


ALTER TABLE public.scan_tool_metadata OWNER TO postgres;

--
-- Name: scan_tool_step; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.scan_tool_step (
    id integer DEFAULT nextval('public.id_seq_scan_tool_step'::regclass) NOT NULL,
    scan_tool_id integer,
    index integer,
    step_execution_type character varying(10),
    step_execution_sync boolean NOT NULL,
    retry_count integer,
    execute_step_on_fail integer,
    execute_step_on_pass integer,
    render_input_data_from_step integer,
    http_input_payload jsonb,
    http_method_type character varying(10),
    http_req_headers jsonb,
    http_query_params jsonb,
    cli_command text,
    cli_output_type character varying(10),
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.scan_tool_step OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: script_path_arg_port_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.script_path_arg_port_mapping (
    id integer DEFAULT nextval('public.id_seq_script_path_arg_port_mapping'::regclass) NOT NULL,
    type_of_mapping character varying(255),
    file_path_on_disk text,
    file_path_on_container text,
    command text,
    args text[],
    port_on_local integer,
    port_on_container integer,
    script_id integer,
    deleted boolean,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.script_path_arg_port_mapping OWNER TO postgres;

--
-- Name: self_registration_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.self_registration_roles (
    role character varying(255) NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.self_registration_roles OWNER TO postgres;

--
-- Name: server_action_audit_log; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.server_action_audit_log (
    id integer DEFAULT nextval('public.id_seq_server_action_audit_log'::regclass) NOT NULL,
    action character varying(255) NOT NULL,
    version character varying(255),
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL
);


ALTER TABLE public.server_action_audit_log OWNER TO postgres;

--
-- Name: ses_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ses_config (
    id integer NOT NULL,
    region character varying(50) NOT NULL,
    access_key character varying(250) NOT NULL,
    secret_access_key character varying(250) NOT NULL,
    session_token character varying(250),
    from_email character varying(250) NOT NULL,
    config_name character varying(250),
    description character varying(500),
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    owner_id integer,
    "default" boolean,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.ses_config OWNER TO postgres;

--
-- Name: ses_config_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ses_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ses_config_id_seq OWNER TO postgres;

--
-- Name: ses_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ses_config_id_seq OWNED BY public.ses_config.id;


--
-- Name: slack_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.slack_config (
    id integer NOT NULL,
    web_hook_url character varying(250) NOT NULL,
    config_name character varying(250) NOT NULL,
    description character varying(500),
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer,
    owner_id integer,
    team_id integer,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.slack_config OWNER TO postgres;

--
-- Name: slack_config_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.slack_config_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.slack_config_id_seq OWNER TO postgres;

--
-- Name: slack_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.slack_config_id_seq OWNED BY public.slack_config.id;


--
-- Name: smtp_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.smtp_config (
    id integer DEFAULT nextval('public.id_seq_smtp_config'::regclass) NOT NULL,
    port text,
    host text,
    auth_type text,
    auth_user text,
    auth_password text,
    from_email text,
    config_name text,
    description text,
    owner_id integer,
    "default" boolean,
    deleted boolean DEFAULT false NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.smtp_config OWNER TO postgres;

--
-- Name: sso_login_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sso_login_config (
    id integer DEFAULT nextval('public.id_seq_sso_login_config'::regclass) NOT NULL,
    name character varying(250),
    label character varying(250),
    url character varying(250),
    config text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer,
    active boolean
);


ALTER TABLE public.sso_login_config OWNER TO postgres;

--
-- Name: system_network_controller_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.system_network_controller_config (
    id integer DEFAULT nextval('public.id_seq_system_network_controller_config'::regclass) NOT NULL,
    ip character varying(50),
    username character varying(100),
    password character varying(100),
    active boolean,
    action_link_json jsonb,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.system_network_controller_config OWNER TO postgres;

--
-- Name: team; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.team (
    id integer NOT NULL,
    name character varying(250) NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.team OWNER TO postgres;

--
-- Name: team_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.team_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.team_id_seq OWNER TO postgres;

--
-- Name: team_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.team_id_seq OWNED BY public.team.id;


--
-- Name: terminal_access_templates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.terminal_access_templates (
    id integer DEFAULT nextval('public.id_seq_terminal_access_templates'::regclass) NOT NULL,
    template_name character varying(1000),
    template_data text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.terminal_access_templates OWNER TO postgres;

--
-- Name: timeout_window_configuration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.timeout_window_configuration (
    id integer DEFAULT nextval('public.id_seq_timeout_window_configuration'::regclass) NOT NULL,
    timeout_window_expression character varying(255) NOT NULL,
    timeout_window_expression_format integer NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.timeout_window_configuration OWNER TO postgres;

--
-- Name: timeout_window_resource_mappings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.timeout_window_resource_mappings (
    id integer DEFAULT nextval('public.id_seq_timeout_window_resource_mappings'::regclass) NOT NULL,
    timeout_window_configuration_id integer NOT NULL,
    resource_id integer NOT NULL,
    resource_type integer NOT NULL
);


ALTER TABLE public.timeout_window_resource_mappings OWNER TO postgres;

--
-- Name: user_attributes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_attributes (
    email_id character varying(500) NOT NULL,
    user_data json NOT NULL,
    created_on timestamp with time zone,
    updated_on timestamp with time zone,
    created_by integer,
    updated_by integer
);


ALTER TABLE public.user_attributes OWNER TO postgres;

--
-- Name: user_audit; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_audit (
    id integer DEFAULT nextval('public.id_seq_user_audit'::regclass) NOT NULL,
    user_id integer NOT NULL,
    client_ip character varying(256) NOT NULL,
    created_on timestamp with time zone NOT NULL,
    updated_on timestamp with time zone
);


ALTER TABLE public.user_audit OWNER TO postgres;

--
-- Name: user_auto_assigned_groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_auto_assigned_groups (
    id integer DEFAULT nextval('public.id_seq_user_groups'::regclass) NOT NULL,
    user_id integer,
    group_name text,
    is_group_claims_data boolean NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.user_auto_assigned_groups OWNER TO postgres;

--
-- Name: user_deployment_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_deployment_request (
    id integer DEFAULT nextval('public.id_seq_user_deployment_request_sequence'::regclass) NOT NULL,
    pipeline_id integer NOT NULL,
    ci_artifact_id integer NOT NULL,
    additional_override bytea,
    force_trigger boolean DEFAULT false NOT NULL,
    force_sync_deployment boolean DEFAULT false NOT NULL,
    strategy character varying(100),
    deployment_with_config character varying(100),
    specific_trigger_wfr_id integer,
    cd_workflow_id integer NOT NULL,
    deployment_type integer,
    triggered_at timestamp with time zone NOT NULL,
    triggered_by integer NOT NULL
);


ALTER TABLE public.user_deployment_request OWNER TO postgres;

--
-- Name: user_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_group (
    id integer DEFAULT nextval('public.id_seq_user_group'::regclass) NOT NULL,
    name character varying(50) NOT NULL,
    identifier character varying(50) NOT NULL,
    description text NOT NULL,
    active boolean NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.user_group OWNER TO postgres;

--
-- Name: user_group_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_group_mapping (
    id integer DEFAULT nextval('public.id_seq_user_group_mapping'::regclass) NOT NULL,
    user_id integer NOT NULL,
    user_group_id integer NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.user_group_mapping OWNER TO postgres;

--
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_roles_id_seq OWNER TO postgres;

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_roles (
    id integer DEFAULT nextval('public.user_roles_id_seq'::regclass) NOT NULL,
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    created_by integer,
    created_on timestamp without time zone,
    updated_by integer,
    updated_on timestamp without time zone,
    timeout_window_configuration_id integer
);


ALTER TABLE public.user_roles OWNER TO postgres;

--
-- Name: user_terminal_access_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_terminal_access_data (
    id integer DEFAULT nextval('public.id_seq_user_terminal_access_data'::regclass) NOT NULL,
    user_id integer,
    cluster_id integer,
    pod_name character varying(1000),
    node_name character varying(1000),
    status character varying(1000),
    metadata json,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.user_terminal_access_data OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer DEFAULT nextval('public.users_id_seq'::regclass) NOT NULL,
    fname text,
    lname text,
    password text,
    access_token text,
    created_on timestamp without time zone,
    email_id character varying(100) NOT NULL,
    created_by integer,
    updated_by integer,
    updated_on timestamp without time zone,
    active boolean DEFAULT true NOT NULL,
    user_type character varying(250),
    timeout_window_configuration_id integer,
    request_email_id character varying(256)
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: value_constraint; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.value_constraint (
    id integer DEFAULT nextval('public.id_seq_value_constraint'::regclass) NOT NULL,
    choices text[],
    value_of character varying(50) NOT NULL,
    block_custom_value boolean DEFAULT false NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    "constraint" jsonb NOT NULL,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.value_constraint OWNER TO postgres;

--
-- Name: variable_data_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variable_data_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variable_data_seq OWNER TO postgres;

--
-- Name: variable_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variable_data (
    id integer DEFAULT nextval('public.variable_data_seq'::regclass) NOT NULL,
    variable_scope_id integer,
    data text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.variable_data OWNER TO postgres;

--
-- Name: variable_definition_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variable_definition_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variable_definition_seq OWNER TO postgres;

--
-- Name: variable_definition; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variable_definition (
    id integer DEFAULT nextval('public.variable_definition_seq'::regclass) NOT NULL,
    name character varying(300) NOT NULL,
    data_type character varying(50) NOT NULL,
    var_type character varying(50) NOT NULL,
    active boolean NOT NULL,
    description text,
    short_description text,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.variable_definition OWNER TO postgres;

--
-- Name: variable_entity_mapping_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variable_entity_mapping_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variable_entity_mapping_seq OWNER TO postgres;

--
-- Name: variable_entity_mapping; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variable_entity_mapping (
    id integer DEFAULT nextval('public.variable_entity_mapping_seq'::regclass) NOT NULL,
    variable_name character varying(300) NOT NULL,
    entity_type integer NOT NULL,
    entity_id integer NOT NULL,
    is_deleted boolean NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.variable_entity_mapping OWNER TO postgres;

--
-- Name: variable_snapshot_history_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.variable_snapshot_history_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.variable_snapshot_history_seq OWNER TO postgres;

--
-- Name: variable_snapshot_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.variable_snapshot_history (
    id integer DEFAULT nextval('public.variable_snapshot_history_seq'::regclass) NOT NULL,
    variable_snapshot jsonb NOT NULL,
    history_reference_type integer NOT NULL,
    history_reference_id integer NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.variable_snapshot_history OWNER TO postgres;

--
-- Name: webhook_config; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.webhook_config (
    id integer DEFAULT nextval('public.id_seq_webhook_config'::regclass) NOT NULL,
    web_hook_url character varying,
    config_name character varying(100),
    header jsonb,
    payload character varying,
    description text,
    owner_id integer,
    active boolean,
    deleted boolean DEFAULT false NOT NULL,
    created_on timestamp with time zone,
    created_by integer,
    updated_on timestamp with time zone,
    updated_by integer
);


ALTER TABLE public.webhook_config OWNER TO postgres;

--
-- Name: webhook_event_data_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.webhook_event_data_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.webhook_event_data_id_seq OWNER TO postgres;

--
-- Name: webhook_event_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.webhook_event_data (
    id integer DEFAULT nextval('public.webhook_event_data_id_seq'::regclass) NOT NULL,
    git_host_id integer NOT NULL,
    event_type character varying(250) NOT NULL,
    payload_json json NOT NULL,
    created_on timestamp with time zone NOT NULL
);


ALTER TABLE public.webhook_event_data OWNER TO postgres;

--
-- Name: workflow_execution_stage; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.workflow_execution_stage (
    id integer DEFAULT nextval('public.id_seq_workflow_execution_stage'::regclass) NOT NULL,
    stage_name character varying(50),
    step_name character varying(50),
    status character varying(50),
    status_for character varying(50),
    message text,
    metadata text,
    workflow_id integer NOT NULL,
    workflow_type character varying(50) NOT NULL,
    start_time text,
    end_time text,
    created_on timestamp with time zone NOT NULL,
    created_by integer NOT NULL,
    updated_on timestamp with time zone NOT NULL,
    updated_by integer NOT NULL
);


ALTER TABLE public.workflow_execution_stage OWNER TO postgres;

--
-- Name: app id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app ALTER COLUMN id SET DEFAULT nextval('public.app_id_seq'::regclass);


--
-- Name: app_env_linkouts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_env_linkouts ALTER COLUMN id SET DEFAULT nextval('public.app_env_linkouts_id_seq'::regclass);


--
-- Name: app_level_metrics id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_level_metrics ALTER COLUMN id SET DEFAULT nextval('public.app_level_metrics_id_seq'::regclass);


--
-- Name: app_store id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store ALTER COLUMN id SET DEFAULT nextval('public.app_store_id_seq'::regclass);


--
-- Name: app_store_application_version id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store_application_version ALTER COLUMN id SET DEFAULT nextval('public.app_store_application_version_id_seq'::regclass);


--
-- Name: app_store_version_values id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store_version_values ALTER COLUMN id SET DEFAULT nextval('public.app_store_version_values_id_seq'::regclass);


--
-- Name: app_workflow id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_workflow ALTER COLUMN id SET DEFAULT nextval('public.app_workflow_id_seq'::regclass);


--
-- Name: cd_workflow id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow ALTER COLUMN id SET DEFAULT nextval('public.cd_workflow_id_seq'::regclass);


--
-- Name: cd_workflow_runner id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow_runner ALTER COLUMN id SET DEFAULT nextval('public.cd_workflow_runner_id_seq'::regclass);


--
-- Name: chart_env_config_override id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_env_config_override ALTER COLUMN id SET DEFAULT nextval('public.chart_env_config_override_id_seq'::regclass);


--
-- Name: chart_group id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group ALTER COLUMN id SET DEFAULT nextval('public.chart_group_id_seq'::regclass);


--
-- Name: chart_group_deployment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_deployment ALTER COLUMN id SET DEFAULT nextval('public.chart_group_deployment_id_seq'::regclass);


--
-- Name: chart_group_entry id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_entry ALTER COLUMN id SET DEFAULT nextval('public.chart_group_entry_id_seq'::regclass);


--
-- Name: chart_ref_schema id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_ref_schema ALTER COLUMN id SET DEFAULT nextval('public.chart_ref_schema_id_seq'::regclass);


--
-- Name: chart_repo id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_repo ALTER COLUMN id SET DEFAULT nextval('public.chart_repo_id_seq'::regclass);


--
-- Name: charts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charts ALTER COLUMN id SET DEFAULT nextval('public.charts_id_seq'::regclass);


--
-- Name: ci_artifact id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact ALTER COLUMN id SET DEFAULT nextval('public.ci_artifact_id_seq'::regclass);


--
-- Name: ci_pipeline id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline ALTER COLUMN id SET DEFAULT nextval('public.ci_pipeline_id_seq'::regclass);


--
-- Name: ci_pipeline_material id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_material ALTER COLUMN id SET DEFAULT nextval('public.ci_pipeline_material_id_seq'::regclass);


--
-- Name: ci_pipeline_scripts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_scripts ALTER COLUMN id SET DEFAULT nextval('public.ci_pipeline_scripts_id_seq'::regclass);


--
-- Name: ci_template id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template ALTER COLUMN id SET DEFAULT nextval('public.ci_template_id_seq'::regclass);


--
-- Name: ci_workflow id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_workflow ALTER COLUMN id SET DEFAULT nextval('public.ci_workflow_id_seq'::regclass);


--
-- Name: cluster id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster ALTER COLUMN id SET DEFAULT nextval('public.cluster_id_seq'::regclass);


--
-- Name: cluster_accounts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_accounts ALTER COLUMN id SET DEFAULT nextval('public.cluster_accounts_id_seq'::regclass);


--
-- Name: cluster_helm_config id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_helm_config ALTER COLUMN id SET DEFAULT nextval('public.cluster_helm_config_id_seq'::regclass);


--
-- Name: custom_tag id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.custom_tag ALTER COLUMN id SET DEFAULT nextval('public.custom_tag_id_seq'::regclass);


--
-- Name: deployment_group id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group ALTER COLUMN id SET DEFAULT nextval('public.deployment_group_id_seq'::regclass);


--
-- Name: deployment_group_app id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group_app ALTER COLUMN id SET DEFAULT nextval('public.deployment_group_app_id_seq'::regclass);


--
-- Name: deployment_status id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_status ALTER COLUMN id SET DEFAULT nextval('public.deployment_status_id_seq'::regclass);


--
-- Name: env_level_app_metrics id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.env_level_app_metrics ALTER COLUMN id SET DEFAULT nextval('public.env_level_app_metrics_id_seq'::regclass);


--
-- Name: environment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.environment ALTER COLUMN id SET DEFAULT nextval('public.environment_id_seq'::regclass);


--
-- Name: event id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event ALTER COLUMN id SET DEFAULT nextval('public.event_id_seq'::regclass);


--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: external_ci_pipeline id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_ci_pipeline ALTER COLUMN id SET DEFAULT nextval('public.external_ci_pipeline_id_seq'::regclass);


--
-- Name: git_material id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_material ALTER COLUMN id SET DEFAULT nextval('public.git_material_id_seq'::regclass);


--
-- Name: git_provider id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_provider ALTER COLUMN id SET DEFAULT nextval('public.git_provider_id_seq'::regclass);


--
-- Name: git_web_hook id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_web_hook ALTER COLUMN id SET DEFAULT nextval('public.git_web_hook_id_seq'::regclass);


--
-- Name: image_path_reservation id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_path_reservation ALTER COLUMN id SET DEFAULT nextval('public.image_path_reservation_id_seq'::regclass);


--
-- Name: installed_app_versions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_app_versions ALTER COLUMN id SET DEFAULT nextval('public.installed_app_versions_id_seq'::regclass);


--
-- Name: installed_apps id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_apps ALTER COLUMN id SET DEFAULT nextval('public.installed_apps_id_seq'::regclass);


--
-- Name: job_event id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_event ALTER COLUMN id SET DEFAULT nextval('public.job_event_id_seq'::regclass);


--
-- Name: notification_settings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings ALTER COLUMN id SET DEFAULT nextval('public.notification_settings_id_seq'::regclass);


--
-- Name: notification_settings_view id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings_view ALTER COLUMN id SET DEFAULT nextval('public.notification_settings_view_id_seq'::regclass);


--
-- Name: notification_templates id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_templates ALTER COLUMN id SET DEFAULT nextval('public.notification_templates_id_seq'::regclass);


--
-- Name: notifier_event_log id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifier_event_log ALTER COLUMN id SET DEFAULT nextval('public.notifier_event_log_id_seq'::regclass);


--
-- Name: pipeline id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline ALTER COLUMN id SET DEFAULT nextval('public.pipeline_id_seq'::regclass);


--
-- Name: pipeline_config_override id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_config_override ALTER COLUMN id SET DEFAULT nextval('public.pipeline_config_override_id_seq'::regclass);


--
-- Name: pipeline_strategy id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_strategy ALTER COLUMN id SET DEFAULT nextval('public.pipeline_strategy_id_seq'::regclass);


--
-- Name: role_group id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_group ALTER COLUMN id SET DEFAULT nextval('public.role_group_id_seq'::regclass);


--
-- Name: role_group_role_mapping id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_group_role_mapping ALTER COLUMN id SET DEFAULT nextval('public.role_group_role_mapping_id_seq'::regclass);


--
-- Name: ses_config id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ses_config ALTER COLUMN id SET DEFAULT nextval('public.ses_config_id_seq'::regclass);


--
-- Name: slack_config id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.slack_config ALTER COLUMN id SET DEFAULT nextval('public.slack_config_id_seq'::regclass);


--
-- Name: team id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.team ALTER COLUMN id SET DEFAULT nextval('public.team_id_seq'::regclass);


--
-- Name: api_token api_token_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_token
    ADD CONSTRAINT api_token_name_key UNIQUE (name);


--
-- Name: api_token api_token_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_token
    ADD CONSTRAINT api_token_pkey PRIMARY KEY (id);


--
-- Name: api_token api_token_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_token
    ADD CONSTRAINT api_token_token_key UNIQUE (token);


--
-- Name: app_env_linkouts app_env_linkouts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_env_linkouts
    ADD CONSTRAINT app_env_linkouts_pkey PRIMARY KEY (id);


--
-- Name: resource_group_mapping app_group_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_group_mapping
    ADD CONSTRAINT app_group_mapping_pkey PRIMARY KEY (id);


--
-- Name: resource_group app_group_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_group
    ADD CONSTRAINT app_group_pkey PRIMARY KEY (id);


--
-- Name: app_label app_label_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_label
    ADD CONSTRAINT app_label_pkey PRIMARY KEY (id);


--
-- Name: app_level_metrics app_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_level_metrics
    ADD CONSTRAINT app_metrics_pkey PRIMARY KEY (id);


--
-- Name: app app_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app
    ADD CONSTRAINT app_pkey PRIMARY KEY (id);


--
-- Name: app_status app_status_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_status
    ADD CONSTRAINT app_status_pkey PRIMARY KEY (app_id, env_id);


--
-- Name: app_store_application_version app_store_application_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store_application_version
    ADD CONSTRAINT app_store_application_version_pkey PRIMARY KEY (id);


--
-- Name: app_store app_store_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store
    ADD CONSTRAINT app_store_pkey PRIMARY KEY (id);


--
-- Name: app_store_version_values app_store_version_values_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store_version_values
    ADD CONSTRAINT app_store_version_values_pkey PRIMARY KEY (id);


--
-- Name: app_workflow_mapping app_workflow_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_workflow_mapping
    ADD CONSTRAINT app_workflow_mapping_pkey PRIMARY KEY (id);


--
-- Name: app_workflow app_workflow_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_workflow
    ADD CONSTRAINT app_workflow_pkey PRIMARY KEY (id);


--
-- Name: artifact_promotion_approval_request artifact_promotion_approval_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.artifact_promotion_approval_request
    ADD CONSTRAINT artifact_promotion_approval_request_pkey PRIMARY KEY (id);


--
-- Name: attributes attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);


--
-- Name: auto_remediation_trigger auto_remediation_trigger_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auto_remediation_trigger
    ADD CONSTRAINT auto_remediation_trigger_pkey PRIMARY KEY (id);


--
-- Name: bulk_update_readme bulk_update_readme_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bulk_update_readme
    ADD CONSTRAINT bulk_update_readme_pkey PRIMARY KEY (id);


--
-- Name: cd_workflow cd_workflow_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow
    ADD CONSTRAINT cd_workflow_pkey PRIMARY KEY (id);


--
-- Name: cd_workflow_runner cd_workflow_runner_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow_runner
    ADD CONSTRAINT cd_workflow_runner_pkey PRIMARY KEY (id);


--
-- Name: chart_category_mapping chart_category_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_category_mapping
    ADD CONSTRAINT chart_category_mapping_pkey PRIMARY KEY (id);


--
-- Name: chart_category chart_category_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_category
    ADD CONSTRAINT chart_category_pkey PRIMARY KEY (id);


--
-- Name: chart_env_config_override chart_env_config_override_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_env_config_override
    ADD CONSTRAINT chart_env_config_override_pkey PRIMARY KEY (id);


--
-- Name: chart_group_deployment chart_group_deployment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_deployment
    ADD CONSTRAINT chart_group_deployment_pkey PRIMARY KEY (id);


--
-- Name: chart_group_entry chart_group_entry_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_entry
    ADD CONSTRAINT chart_group_entry_pkey PRIMARY KEY (id);


--
-- Name: chart_group chart_group_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group
    ADD CONSTRAINT chart_group_pkey PRIMARY KEY (id);


--
-- Name: chart_ref_metadata chart_ref_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_ref_metadata
    ADD CONSTRAINT chart_ref_metadata_pkey PRIMARY KEY (chart_name);


--
-- Name: chart_ref_schema chart_ref_schema_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_ref_schema
    ADD CONSTRAINT chart_ref_schema_pkey PRIMARY KEY (id);


--
-- Name: chart_repo chart_repo_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_repo
    ADD CONSTRAINT chart_repo_pkey PRIMARY KEY (id);


--
-- Name: charts charts_chart_name_chart_version_chart_repo_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charts
    ADD CONSTRAINT charts_chart_name_chart_version_chart_repo_key UNIQUE (chart_name, chart_version, chart_repo);


--
-- Name: charts charts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charts
    ADD CONSTRAINT charts_pkey PRIMARY KEY (id);


--
-- Name: ci_artifact ci_artifact_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact
    ADD CONSTRAINT ci_artifact_pkey PRIMARY KEY (id);


--
-- Name: ci_build_config ci_build_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_build_config
    ADD CONSTRAINT ci_build_config_pkey PRIMARY KEY (id);


--
-- Name: ci_env_mapping_history ci_env_mapping_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_env_mapping_history
    ADD CONSTRAINT ci_env_mapping_history_pkey PRIMARY KEY (id);


--
-- Name: ci_env_mapping ci_env_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_env_mapping
    ADD CONSTRAINT ci_env_mapping_pkey PRIMARY KEY (id);


--
-- Name: ci_pipeline_history ci_pipeline_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_history
    ADD CONSTRAINT ci_pipeline_history_pkey PRIMARY KEY (id);


--
-- Name: ci_pipeline_material ci_pipeline_material_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_material
    ADD CONSTRAINT ci_pipeline_material_pkey PRIMARY KEY (id);


--
-- Name: ci_pipeline ci_pipeline_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline
    ADD CONSTRAINT ci_pipeline_pkey PRIMARY KEY (id);


--
-- Name: ci_pipeline_scripts ci_pipeline_scripts_name_ci_pipeline_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_scripts
    ADD CONSTRAINT ci_pipeline_scripts_name_ci_pipeline_id_key UNIQUE (name, ci_pipeline_id);


--
-- Name: ci_pipeline_scripts ci_pipeline_scripts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_scripts
    ADD CONSTRAINT ci_pipeline_scripts_pkey PRIMARY KEY (id);


--
-- Name: ci_template ci_template_app_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_app_id_key UNIQUE (app_id);


--
-- Name: ci_template_history ci_template_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_history
    ADD CONSTRAINT ci_template_history_pkey PRIMARY KEY (id);


--
-- Name: ci_template_override ci_template_override_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_override
    ADD CONSTRAINT ci_template_override_pkey PRIMARY KEY (id);


--
-- Name: ci_template ci_template_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_pkey PRIMARY KEY (id);


--
-- Name: ci_workflow ci_workflow_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_workflow
    ADD CONSTRAINT ci_workflow_pkey PRIMARY KEY (id);


--
-- Name: cluster_accounts cluster_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_accounts
    ADD CONSTRAINT cluster_accounts_pkey PRIMARY KEY (id);


--
-- Name: cluster_helm_config cluster_helm_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_helm_config
    ADD CONSTRAINT cluster_helm_config_pkey PRIMARY KEY (id);


--
-- Name: cluster_installed_apps cluster_installed_apps_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_installed_apps
    ADD CONSTRAINT cluster_installed_apps_pkey PRIMARY KEY (id);


--
-- Name: generic_note_history cluster_note_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.generic_note_history
    ADD CONSTRAINT cluster_note_history_pkey PRIMARY KEY (id);


--
-- Name: generic_note cluster_note_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.generic_note
    ADD CONSTRAINT cluster_note_pkey PRIMARY KEY (id);


--
-- Name: cluster cluster_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster
    ADD CONSTRAINT cluster_pkey PRIMARY KEY (id);


--
-- Name: config_map_history config_map_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_history
    ADD CONSTRAINT config_map_history_pkey PRIMARY KEY (id);


--
-- Name: custom_tag custom_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.custom_tag
    ADD CONSTRAINT custom_tag_pkey PRIMARY KEY (id);


--
-- Name: cve_policy_control cve_policy_control_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cve_policy_control
    ADD CONSTRAINT cve_policy_control_pkey PRIMARY KEY (id);


--
-- Name: cve_store cve_store_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cve_store
    ADD CONSTRAINT cve_store_pkey PRIMARY KEY (name);


--
-- Name: default_auth_policy default_auth_policy_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.default_auth_policy
    ADD CONSTRAINT default_auth_policy_pkey PRIMARY KEY (id);


--
-- Name: default_auth_role default_auth_role_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.default_auth_role
    ADD CONSTRAINT default_auth_role_pkey PRIMARY KEY (id);


--
-- Name: default_rbac_role_data default_rbac_role_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.default_rbac_role_data
    ADD CONSTRAINT default_rbac_role_data_pkey PRIMARY KEY (id);


--
-- Name: deployment_app_migration_history deployment_app_migration_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_app_migration_history
    ADD CONSTRAINT deployment_app_migration_history_pkey PRIMARY KEY (id);


--
-- Name: deployment_approval_request deployment_approval_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_approval_request
    ADD CONSTRAINT deployment_approval_request_pkey PRIMARY KEY (id);


--
-- Name: request_approval_user_data deployment_approval_user_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.request_approval_user_data
    ADD CONSTRAINT deployment_approval_user_data_pkey PRIMARY KEY (id);


--
-- Name: deployment_config deployment_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_config
    ADD CONSTRAINT deployment_config_pkey PRIMARY KEY (id);


--
-- Name: deployment_event deployment_event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_event
    ADD CONSTRAINT deployment_event_pkey PRIMARY KEY (id);


--
-- Name: deployment_group_app deployment_group_app_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group_app
    ADD CONSTRAINT deployment_group_app_pkey PRIMARY KEY (id);


--
-- Name: deployment_group deployment_group_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group
    ADD CONSTRAINT deployment_group_pkey PRIMARY KEY (id);


--
-- Name: deployment_template_history deployment_template_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_template_history
    ADD CONSTRAINT deployment_template_history_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_object_audit devtron_resource_object_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_object_audit
    ADD CONSTRAINT devtron_resource_object_audit_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_object_dep_relations devtron_resource_object_dep_relations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_object_dep_relations
    ADD CONSTRAINT devtron_resource_object_dep_relations_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_object devtron_resource_object_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_object
    ADD CONSTRAINT devtron_resource_object_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource devtron_resource_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource
    ADD CONSTRAINT devtron_resource_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_schema_audit devtron_resource_schema_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_schema_audit
    ADD CONSTRAINT devtron_resource_schema_audit_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_schema devtron_resource_schema_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_schema
    ADD CONSTRAINT devtron_resource_schema_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_searchable_key devtron_resource_searchable_key_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_searchable_key
    ADD CONSTRAINT devtron_resource_searchable_key_pkey PRIMARY KEY (id);


--
-- Name: devtron_resource_task_run devtron_resource_task_run_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_task_run
    ADD CONSTRAINT devtron_resource_task_run_pkey PRIMARY KEY (id);


--
-- Name: docker_artifact_store docker_artifact_store_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.docker_artifact_store
    ADD CONSTRAINT docker_artifact_store_pkey PRIMARY KEY (id);


--
-- Name: docker_registry_ips_config docker_registry_ips_config_docker_artifact_store_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.docker_registry_ips_config
    ADD CONSTRAINT docker_registry_ips_config_docker_artifact_store_id_key UNIQUE (docker_artifact_store_id);


--
-- Name: docker_registry_ips_config docker_registry_ips_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.docker_registry_ips_config
    ADD CONSTRAINT docker_registry_ips_config_pkey PRIMARY KEY (id);


--
-- Name: draft draft_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.draft
    ADD CONSTRAINT draft_pkey PRIMARY KEY (id);


--
-- Name: draft_version_comment draft_version_comment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.draft_version_comment
    ADD CONSTRAINT draft_version_comment_pkey PRIMARY KEY (id);


--
-- Name: draft_version draft_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.draft_version
    ADD CONSTRAINT draft_version_pkey PRIMARY KEY (id);


--
-- Name: deployment_status ds_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_status
    ADD CONSTRAINT ds_pkey PRIMARY KEY (id);


--
-- Name: env_level_app_metrics env_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.env_level_app_metrics
    ADD CONSTRAINT env_metrics_pkey PRIMARY KEY (id);


--
-- Name: environment environment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.environment
    ADD CONSTRAINT environment_pkey PRIMARY KEY (id);


--
-- Name: ephemeral_container_actions ephemeral_container_actions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ephemeral_container_actions
    ADD CONSTRAINT ephemeral_container_actions_pkey PRIMARY KEY (id);


--
-- Name: ephemeral_container ephemeral_container_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ephemeral_container
    ADD CONSTRAINT ephemeral_container_pkey PRIMARY KEY (id);


--
-- Name: event event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT event_pkey PRIMARY KEY (id);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: external_ci_pipeline external_ci_pipeline_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_ci_pipeline
    ADD CONSTRAINT external_ci_pipeline_pkey PRIMARY KEY (id);


--
-- Name: external_link_identifier_mapping external_link_cluster_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_link_identifier_mapping
    ADD CONSTRAINT external_link_cluster_mapping_pkey PRIMARY KEY (id);


--
-- Name: external_link_monitoring_tool external_link_monitoring_tool_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_link_monitoring_tool
    ADD CONSTRAINT external_link_monitoring_tool_pkey PRIMARY KEY (id);


--
-- Name: external_link external_link_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_link
    ADD CONSTRAINT external_link_pkey PRIMARY KEY (id);


--
-- Name: file_reference file_reference_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.file_reference
    ADD CONSTRAINT file_reference_pkey PRIMARY KEY (id);


--
-- Name: git_host git_host_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_host
    ADD CONSTRAINT git_host_name_key UNIQUE (name);


--
-- Name: git_host git_host_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_host
    ADD CONSTRAINT git_host_pkey PRIMARY KEY (id);


--
-- Name: git_material_history git_material_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_material_history
    ADD CONSTRAINT git_material_history_pkey PRIMARY KEY (id);


--
-- Name: git_material git_material_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_material
    ADD CONSTRAINT git_material_pkey PRIMARY KEY (id);


--
-- Name: git_provider git_provider_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_provider
    ADD CONSTRAINT git_provider_pkey PRIMARY KEY (id);


--
-- Name: git_sensor_node_mapping git_sensor_node_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_sensor_node_mapping
    ADD CONSTRAINT git_sensor_node_mapping_pkey PRIMARY KEY (id);


--
-- Name: git_sensor_node git_sensor_node_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_sensor_node
    ADD CONSTRAINT git_sensor_node_pkey PRIMARY KEY (id);


--
-- Name: git_web_hook git_web_hook_ci_material_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_web_hook
    ADD CONSTRAINT git_web_hook_ci_material_id_key UNIQUE (ci_material_id);


--
-- Name: git_web_hook git_web_hook_git_material_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_web_hook
    ADD CONSTRAINT git_web_hook_git_material_id_key UNIQUE (git_material_id);


--
-- Name: git_web_hook git_web_hook_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_web_hook
    ADD CONSTRAINT git_web_hook_pkey PRIMARY KEY (id);


--
-- Name: gitops_config gitops_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gitops_config
    ADD CONSTRAINT gitops_config_pkey PRIMARY KEY (id);


--
-- Name: global_authorisation_config global_authorisation_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_authorisation_config
    ADD CONSTRAINT global_authorisation_config_pkey PRIMARY KEY (id);


--
-- Name: global_cm_cs global_cm_cs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_cm_cs
    ADD CONSTRAINT global_cm_cs_pkey PRIMARY KEY (id);


--
-- Name: global_policy_history global_policy_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_policy_history
    ADD CONSTRAINT global_policy_history_pkey PRIMARY KEY (id);


--
-- Name: global_policy global_policy_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_policy
    ADD CONSTRAINT global_policy_pkey PRIMARY KEY (id);


--
-- Name: global_policy_searchable_field global_policy_searchable_field_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_policy_searchable_field
    ADD CONSTRAINT global_policy_searchable_field_pkey PRIMARY KEY (id);


--
-- Name: global_strategy_metadata_chart_ref_mapping global_strategy_metadata_chart_ref_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_strategy_metadata_chart_ref_mapping
    ADD CONSTRAINT global_strategy_metadata_chart_ref_mapping_pkey PRIMARY KEY (id);


--
-- Name: global_strategy_metadata global_strategy_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_strategy_metadata
    ADD CONSTRAINT global_strategy_metadata_pkey PRIMARY KEY (id);


--
-- Name: global_tag global_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_tag
    ADD CONSTRAINT global_tag_pkey PRIMARY KEY (id);


--
-- Name: helm_values helm_values_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.helm_values
    ADD CONSTRAINT helm_values_pkey PRIMARY KEY (app_name, environment);


--
-- Name: image_comments image_comments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_comments
    ADD CONSTRAINT image_comments_pkey PRIMARY KEY (id);


--
-- Name: image_path_reservation image_path_reservation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_path_reservation
    ADD CONSTRAINT image_path_reservation_pkey PRIMARY KEY (id);


--
-- Name: image_scan_deploy_info image_scan_deploy_info_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_deploy_info
    ADD CONSTRAINT image_scan_deploy_info_pkey PRIMARY KEY (id);


--
-- Name: image_scan_execution_history image_scan_execution_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_execution_history
    ADD CONSTRAINT image_scan_execution_history_pkey PRIMARY KEY (id);


--
-- Name: image_scan_execution_result image_scan_execution_result_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_execution_result
    ADD CONSTRAINT image_scan_execution_result_pkey PRIMARY KEY (id);


--
-- Name: image_scan_object_meta image_scan_object_meta_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_object_meta
    ADD CONSTRAINT image_scan_object_meta_pkey PRIMARY KEY (id);


--
-- Name: image_tagging_audit image_tagging_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_tagging_audit
    ADD CONSTRAINT image_tagging_audit_pkey PRIMARY KEY (id);


--
-- Name: infra_config_trigger_history infra_config_trigger_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infra_config_trigger_history
    ADD CONSTRAINT infra_config_trigger_history_pkey PRIMARY KEY (id);


--
-- Name: infra_config_trigger_history infra_config_trigger_history_workflow_id_workflow_type_key__key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infra_config_trigger_history
    ADD CONSTRAINT infra_config_trigger_history_workflow_id_workflow_type_key__key UNIQUE (workflow_id, workflow_type, key, platform);


--
-- Name: infra_profile_configuration infra_profile_configuration_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infra_profile_configuration
    ADD CONSTRAINT infra_profile_configuration_pkey PRIMARY KEY (id);


--
-- Name: infra_profile infra_profile_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infra_profile
    ADD CONSTRAINT infra_profile_pkey PRIMARY KEY (id);


--
-- Name: infrastructure_installation infrastructure_installation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infrastructure_installation
    ADD CONSTRAINT infrastructure_installation_pkey PRIMARY KEY (id);


--
-- Name: infrastructure_installation_versions infrastructure_installation_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infrastructure_installation_versions
    ADD CONSTRAINT infrastructure_installation_versions_pkey PRIMARY KEY (id);


--
-- Name: installed_app_version_history installed_app_version_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_app_version_history
    ADD CONSTRAINT installed_app_version_history_pkey PRIMARY KEY (id);


--
-- Name: installed_app_versions installed_app_versions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_app_versions
    ADD CONSTRAINT installed_app_versions_pkey PRIMARY KEY (id);


--
-- Name: installed_apps installed_apps_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_apps
    ADD CONSTRAINT installed_apps_pkey PRIMARY KEY (id);


--
-- Name: intercepted_event_execution intercepted_event_execution_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercepted_event_execution
    ADD CONSTRAINT intercepted_event_execution_pkey PRIMARY KEY (id);


--
-- Name: job_event job_event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.job_event
    ADD CONSTRAINT job_event_pkey PRIMARY KEY (id);


--
-- Name: k8s_event_watcher k8s_event_watcher_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.k8s_event_watcher
    ADD CONSTRAINT k8s_event_watcher_pkey PRIMARY KEY (id);


--
-- Name: kubernetes_resource_history kubernetes_resource_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.kubernetes_resource_history
    ADD CONSTRAINT kubernetes_resource_history_pkey PRIMARY KEY (id);


--
-- Name: license_attributes license_attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.license_attributes
    ADD CONSTRAINT license_attributes_pkey PRIMARY KEY (id);


--
-- Name: lock_configuration lock_configuration_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.lock_configuration
    ADD CONSTRAINT lock_configuration_pkey PRIMARY KEY (id);


--
-- Name: manifest_push_config manifest_push_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.manifest_push_config
    ADD CONSTRAINT manifest_push_config_pkey PRIMARY KEY (id);


--
-- Name: module_action_audit_log module_action_audit_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.module_action_audit_log
    ADD CONSTRAINT module_action_audit_log_pkey PRIMARY KEY (id);


--
-- Name: module module_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.module
    ADD CONSTRAINT module_name_key UNIQUE (name);


--
-- Name: module module_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.module
    ADD CONSTRAINT module_pkey PRIMARY KEY (id);


--
-- Name: module_resource_status module_resource_status_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.module_resource_status
    ADD CONSTRAINT module_resource_status_pkey PRIMARY KEY (id);


--
-- Name: notification_rule notification_rule_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_rule
    ADD CONSTRAINT notification_rule_pkey PRIMARY KEY (id);


--
-- Name: notification_settings notification_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_pkey PRIMARY KEY (id);


--
-- Name: notification_settings_view notification_settings_view_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings_view
    ADD CONSTRAINT notification_settings_view_pkey PRIMARY KEY (id);


--
-- Name: notification_templates notification_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_templates_pkey PRIMARY KEY (id);


--
-- Name: notifier_event_log notifier_event_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifier_event_log
    ADD CONSTRAINT notifier_event_log_pkey PRIMARY KEY (id);


--
-- Name: oci_registry_config oci_registry_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oci_registry_config
    ADD CONSTRAINT oci_registry_config_pkey PRIMARY KEY (id);


--
-- Name: operation_audit operation_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operation_audit
    ADD CONSTRAINT operation_audit_pkey PRIMARY KEY (id);


--
-- Name: panel panel_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.panel
    ADD CONSTRAINT panel_pkey PRIMARY KEY (id);


--
-- Name: pipeline_config_override pipeline_config_override_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_config_override
    ADD CONSTRAINT pipeline_config_override_pkey PRIMARY KEY (id);


--
-- Name: pipeline pipeline_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline
    ADD CONSTRAINT pipeline_pkey PRIMARY KEY (id);


--
-- Name: pipeline_stage pipeline_stage_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage
    ADD CONSTRAINT pipeline_stage_pkey PRIMARY KEY (id);


--
-- Name: pipeline_stage_step_condition pipeline_stage_step_condition_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_condition
    ADD CONSTRAINT pipeline_stage_step_condition_pkey PRIMARY KEY (id);


--
-- Name: pipeline_stage_step pipeline_stage_step_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step
    ADD CONSTRAINT pipeline_stage_step_pkey PRIMARY KEY (id);


--
-- Name: pipeline_stage_step_variable pipeline_stage_step_variable_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_variable
    ADD CONSTRAINT pipeline_stage_step_variable_pkey PRIMARY KEY (id);


--
-- Name: pipeline_status_timeline pipeline_status_timeline_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline
    ADD CONSTRAINT pipeline_status_timeline_pkey PRIMARY KEY (id);


--
-- Name: pipeline_status_timeline_resources pipeline_status_timeline_resources_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline_resources
    ADD CONSTRAINT pipeline_status_timeline_resources_pkey PRIMARY KEY (id);


--
-- Name: pipeline_status_timeline_sync_detail pipeline_status_timeline_sync_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline_sync_detail
    ADD CONSTRAINT pipeline_status_timeline_sync_detail_pkey PRIMARY KEY (id);


--
-- Name: pipeline_strategy_history pipeline_strategy_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_strategy_history
    ADD CONSTRAINT pipeline_strategy_history_pkey PRIMARY KEY (id);


--
-- Name: pipeline_strategy pipeline_strategy_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_strategy
    ADD CONSTRAINT pipeline_strategy_pkey PRIMARY KEY (id);


--
-- Name: plugin_metadata plugin_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_metadata
    ADD CONSTRAINT plugin_metadata_pkey PRIMARY KEY (id);


--
-- Name: plugin_parent_metadata plugin_parent_metadata_identifier_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_parent_metadata
    ADD CONSTRAINT plugin_parent_metadata_identifier_key UNIQUE (identifier);


--
-- Name: plugin_parent_metadata plugin_parent_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_parent_metadata
    ADD CONSTRAINT plugin_parent_metadata_pkey PRIMARY KEY (id);


--
-- Name: plugin_pipeline_script plugin_pipeline_script_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_pipeline_script
    ADD CONSTRAINT plugin_pipeline_script_pkey PRIMARY KEY (id);


--
-- Name: plugin_stage_mapping plugin_stage_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_stage_mapping
    ADD CONSTRAINT plugin_stage_mapping_pkey PRIMARY KEY (id);


--
-- Name: plugin_step_condition plugin_step_condition_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step_condition
    ADD CONSTRAINT plugin_step_condition_pkey PRIMARY KEY (id);


--
-- Name: plugin_step plugin_step_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step
    ADD CONSTRAINT plugin_step_pkey PRIMARY KEY (id);


--
-- Name: plugin_step_variable plugin_step_variable_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step_variable
    ADD CONSTRAINT plugin_step_variable_pkey PRIMARY KEY (id);


--
-- Name: plugin_tag plugin_tag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_tag
    ADD CONSTRAINT plugin_tag_pkey PRIMARY KEY (id);


--
-- Name: plugin_tag_relation plugin_tag_relation_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_tag_relation
    ADD CONSTRAINT plugin_tag_relation_pkey PRIMARY KEY (id);


--
-- Name: pre_post_cd_script_history pre_post_cd_script_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pre_post_cd_script_history
    ADD CONSTRAINT pre_post_cd_script_history_pkey PRIMARY KEY (id);


--
-- Name: pre_post_ci_script_history pre_post_ci_script_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pre_post_ci_script_history
    ADD CONSTRAINT pre_post_ci_script_history_pkey PRIMARY KEY (id);


--
-- Name: profile_platform_mapping profile_platform_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile_platform_mapping
    ADD CONSTRAINT profile_platform_mapping_pkey PRIMARY KEY (id);


--
-- Name: rbac_policy_data rbac_policy_data_entity_access_type_role_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_policy_data
    ADD CONSTRAINT rbac_policy_data_entity_access_type_role_key UNIQUE (entity, access_type, role);


--
-- Name: rbac_policy_data rbac_policy_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_policy_data
    ADD CONSTRAINT rbac_policy_data_pkey PRIMARY KEY (id);


--
-- Name: rbac_policy_resource_detail rbac_policy_resource_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_policy_resource_detail
    ADD CONSTRAINT rbac_policy_resource_detail_pkey PRIMARY KEY (id);


--
-- Name: rbac_policy_resource_detail rbac_policy_resource_detail_resource_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_policy_resource_detail
    ADD CONSTRAINT rbac_policy_resource_detail_resource_key UNIQUE (resource);


--
-- Name: rbac_role_audit rbac_role_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_role_audit
    ADD CONSTRAINT rbac_role_audit_pkey PRIMARY KEY (id);


--
-- Name: rbac_role_data rbac_role_data_entity_access_type_role_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_role_data
    ADD CONSTRAINT rbac_role_data_entity_access_type_role_key UNIQUE (entity, access_type, role);


--
-- Name: rbac_role_data rbac_role_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_role_data
    ADD CONSTRAINT rbac_role_data_pkey PRIMARY KEY (id);


--
-- Name: rbac_role_resource_detail rbac_role_resource_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_role_resource_detail
    ADD CONSTRAINT rbac_role_resource_detail_pkey PRIMARY KEY (id);


--
-- Name: rbac_role_resource_detail rbac_role_resource_detail_resource_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rbac_role_resource_detail
    ADD CONSTRAINT rbac_role_resource_detail_resource_key UNIQUE (resource);


--
-- Name: registry_index_mapping registry_index_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registry_index_mapping
    ADD CONSTRAINT registry_index_mapping_pkey PRIMARY KEY (id);


--
-- Name: release_tags release_tags_app_id_tag_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.release_tags
    ADD CONSTRAINT release_tags_app_id_tag_name_key UNIQUE (app_id, tag_name);


--
-- Name: release_tags release_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.release_tags
    ADD CONSTRAINT release_tags_pkey PRIMARY KEY (id);


--
-- Name: remote_connection_config remote_connection_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.remote_connection_config
    ADD CONSTRAINT remote_connection_config_pkey PRIMARY KEY (id);


--
-- Name: resource_filter_audit resource_filter_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_filter_audit
    ADD CONSTRAINT resource_filter_audit_pkey PRIMARY KEY (id);


--
-- Name: resource_filter_evaluation_audit resource_filter_evaluation_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_filter_evaluation_audit
    ADD CONSTRAINT resource_filter_evaluation_audit_pkey PRIMARY KEY (id);


--
-- Name: resource_filter resource_filter_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_filter
    ADD CONSTRAINT resource_filter_pkey PRIMARY KEY (id);


--
-- Name: resource_protection_history resource_protection_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_protection_history
    ADD CONSTRAINT resource_protection_history_pkey PRIMARY KEY (id);


--
-- Name: resource_protection resource_protection_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_protection
    ADD CONSTRAINT resource_protection_pkey PRIMARY KEY (id);


--
-- Name: resource_qualifier_mapping_criteria resource_qualifier_mapping_criteria_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_qualifier_mapping_criteria
    ADD CONSTRAINT resource_qualifier_mapping_criteria_pkey PRIMARY KEY (id);


--
-- Name: resource_scan_execution_result resource_scan_execution_result_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_scan_execution_result
    ADD CONSTRAINT resource_scan_execution_result_pkey PRIMARY KEY (id);


--
-- Name: role_group role_group_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_group
    ADD CONSTRAINT role_group_pkey PRIMARY KEY (id);


--
-- Name: role_group_role_mapping role_group_role_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_group_role_mapping
    ADD CONSTRAINT role_group_role_mapping_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: scan_step_condition_mapping scan_step_condition_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_step_condition_mapping
    ADD CONSTRAINT scan_step_condition_mapping_pkey PRIMARY KEY (id);


--
-- Name: scan_step_condition scan_step_condition_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_step_condition
    ADD CONSTRAINT scan_step_condition_pkey PRIMARY KEY (id);


--
-- Name: scan_tool_execution_history_mapping scan_tool_execution_history_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_execution_history_mapping
    ADD CONSTRAINT scan_tool_execution_history_mapping_pkey PRIMARY KEY (id);


--
-- Name: scan_tool_metadata scan_tool_metadata_name_version_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_metadata
    ADD CONSTRAINT scan_tool_metadata_name_version_unique UNIQUE (name, version);


--
-- Name: scan_tool_metadata scan_tool_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_metadata
    ADD CONSTRAINT scan_tool_metadata_pkey PRIMARY KEY (id);


--
-- Name: scan_tool_step scan_tool_step_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_step
    ADD CONSTRAINT scan_tool_step_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: script_path_arg_port_mapping script_path_arg_port_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.script_path_arg_port_mapping
    ADD CONSTRAINT script_path_arg_port_mapping_pkey PRIMARY KEY (id);


--
-- Name: server_action_audit_log server_action_audit_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.server_action_audit_log
    ADD CONSTRAINT server_action_audit_log_pkey PRIMARY KEY (id);


--
-- Name: ses_config ses_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ses_config
    ADD CONSTRAINT ses_config_pkey PRIMARY KEY (id);


--
-- Name: slack_config slack_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.slack_config
    ADD CONSTRAINT slack_config_pkey PRIMARY KEY (id);


--
-- Name: smtp_config smtp_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.smtp_config
    ADD CONSTRAINT smtp_config_pkey PRIMARY KEY (id);


--
-- Name: sso_login_config sso_login_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sso_login_config
    ADD CONSTRAINT sso_login_config_pkey PRIMARY KEY (id);


--
-- Name: system_network_controller_config system_network_controller_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.system_network_controller_config
    ADD CONSTRAINT system_network_controller_config_pkey PRIMARY KEY (id);


--
-- Name: team team_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.team
    ADD CONSTRAINT team_pkey PRIMARY KEY (id);


--
-- Name: terminal_access_templates terminal_access_template_name_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terminal_access_templates
    ADD CONSTRAINT terminal_access_template_name_unique UNIQUE (template_name);


--
-- Name: terminal_access_templates terminal_access_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.terminal_access_templates
    ADD CONSTRAINT terminal_access_templates_pkey PRIMARY KEY (id);


--
-- Name: timeout_window_configuration timeout_window_configuration_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timeout_window_configuration
    ADD CONSTRAINT timeout_window_configuration_pkey PRIMARY KEY (id);


--
-- Name: timeout_window_resource_mappings timeout_window_resource_mappings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timeout_window_resource_mappings
    ADD CONSTRAINT timeout_window_resource_mappings_pkey PRIMARY KEY (id);


--
-- Name: ci_artifact unique_ci_workflow_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact
    ADD CONSTRAINT unique_ci_workflow_id UNIQUE (ci_workflow_id);


--
-- Name: custom_tag unique_entity_key_entity_value; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.custom_tag
    ADD CONSTRAINT unique_entity_key_entity_value UNIQUE (entity_key, entity_value);


--
-- Name: event unq_event_name_type; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event
    ADD CONSTRAINT unq_event_name_type UNIQUE (event_type);


--
-- Name: notification_templates unq_notification_template; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT unq_notification_template UNIQUE (channel_type, node_type, event_type_id);


--
-- Name: notification_settings unq_source; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT unq_source UNIQUE (app_id, env_id, pipeline_id, pipeline_type, event_type_id);


--
-- Name: user_attributes user_attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_attributes
    ADD CONSTRAINT user_attributes_pkey PRIMARY KEY (email_id);


--
-- Name: user_audit user_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_audit
    ADD CONSTRAINT user_audit_pkey PRIMARY KEY (id);


--
-- Name: user_deployment_request user_deployment_request_cd_workflow_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_deployment_request
    ADD CONSTRAINT user_deployment_request_cd_workflow_id_key UNIQUE (cd_workflow_id);


--
-- Name: user_deployment_request user_deployment_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_deployment_request
    ADD CONSTRAINT user_deployment_request_pkey PRIMARY KEY (id);


--
-- Name: user_group_mapping user_group_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_group_mapping
    ADD CONSTRAINT user_group_mapping_pkey PRIMARY KEY (id);


--
-- Name: user_group user_group_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_group
    ADD CONSTRAINT user_group_pkey PRIMARY KEY (id);


--
-- Name: user_auto_assigned_groups user_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_auto_assigned_groups
    ADD CONSTRAINT user_groups_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: user_terminal_access_data user_terminal_access_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_terminal_access_data
    ADD CONSTRAINT user_terminal_access_data_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: value_constraint value_constraint_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.value_constraint
    ADD CONSTRAINT value_constraint_pkey PRIMARY KEY (id);


--
-- Name: variable_data variable_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variable_data
    ADD CONSTRAINT variable_data_pkey PRIMARY KEY (id);


--
-- Name: variable_definition variable_definition_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variable_definition
    ADD CONSTRAINT variable_definition_pkey PRIMARY KEY (id);


--
-- Name: variable_entity_mapping variable_entity_mapping_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variable_entity_mapping
    ADD CONSTRAINT variable_entity_mapping_pkey PRIMARY KEY (id);


--
-- Name: resource_qualifier_mapping variable_scope_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_qualifier_mapping
    ADD CONSTRAINT variable_scope_pkey PRIMARY KEY (id);


--
-- Name: variable_snapshot_history variable_snapshot_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.variable_snapshot_history
    ADD CONSTRAINT variable_snapshot_history_pkey PRIMARY KEY (id);


--
-- Name: webhook_config webhook_config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_config
    ADD CONSTRAINT webhook_config_pkey PRIMARY KEY (id);


--
-- Name: webhook_event_data webhook_event_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_event_data
    ADD CONSTRAINT webhook_event_data_pkey PRIMARY KEY (id);


--
-- Name: workflow_execution_stage workflow_execution_stage_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.workflow_execution_stage
    ADD CONSTRAINT workflow_execution_stage_pkey PRIMARY KEY (id);


--
-- Name: app_env_pipeline_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX app_env_pipeline_unique ON public.config_map_pipeline_level USING btree (app_id, environment_id, pipeline_id);


--
-- Name: app_env_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX app_env_unique ON public.config_map_env_level USING btree (app_id, environment_id);


--
-- Name: app_id_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX app_id_unique ON public.config_map_app_level USING btree (app_id);


--
-- Name: app_store_application_version_app_store_id_ix; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX app_store_application_version_app_store_id_ix ON public.app_store_application_version USING btree (app_store_id);


--
-- Name: app_store_unique_chart_repo; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX app_store_unique_chart_repo ON public.app_store USING btree (name, chart_repo_id) WHERE (active = true);


--
-- Name: app_store_unique_oci_repo; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX app_store_unique_oci_repo ON public.app_store USING btree (name, docker_artifact_store_id) WHERE (active = true);


--
-- Name: app_store_version_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX app_store_version_unique ON public.app_store_application_version USING btree (app_store_id, name, version);


--
-- Name: bulk_update_readme_resource_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX bulk_update_readme_resource_idx ON public.bulk_update_readme USING btree (resource);


--
-- Name: cdwf_pipeline_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX cdwf_pipeline_id_idx ON public.cd_workflow USING btree (pipeline_id);


--
-- Name: cdwfr_cd_workflow_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX cdwfr_cd_workflow_id_idx ON public.cd_workflow_runner USING btree (cd_workflow_id);


--
-- Name: chart_ref_schema_unique_active_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX chart_ref_schema_unique_active_idx ON public.chart_ref_schema USING btree (name, resource_type) WHERE (active = true);


--
-- Name: ds_app_name_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ds_app_name_index ON public.deployment_status USING btree (app_name);


--
-- Name: email_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX email_unique ON public.users USING btree (email_id);


--
-- Name: entity_key_value; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX entity_key_value ON public.custom_tag USING btree (entity_key, entity_value);


--
-- Name: events_component; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_component ON public.events USING btree (component);


--
-- Name: events_creation_time_stamp; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_creation_time_stamp ON public.events USING btree (creation_time_stamp);


--
-- Name: events_kind; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_kind ON public.events USING btree (kind);


--
-- Name: events_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_name ON public.events USING btree (name);


--
-- Name: events_namespace; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_namespace ON public.events USING btree (namespace);


--
-- Name: events_reason; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_reason ON public.events USING btree (reason);


--
-- Name: events_resource_revision; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX events_resource_revision ON public.events USING btree (resource_revision);


--
-- Name: gpsf_value_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX gpsf_value_idx ON public.global_policy_searchable_field USING btree (value);


--
-- Name: idx_pipeline_status_timeline_cd_workflow_runner_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pipeline_status_timeline_cd_workflow_runner_id ON public.pipeline_status_timeline USING btree (cd_workflow_runner_id);


--
-- Name: idx_pipeline_status_timeline_installed_app_version_history_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pipeline_status_timeline_installed_app_version_history_id ON public.pipeline_status_timeline USING btree (installed_app_version_history_id);


--
-- Name: idx_run_source_dependency_identifier; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_run_source_dependency_identifier ON public.devtron_resource_task_run USING btree (run_source_dependency_identifier);


--
-- Name: idx_run_source_identifier; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_run_source_identifier ON public.devtron_resource_task_run USING btree (run_source_identifier);


--
-- Name: idx_unique_artifact_promoted_to_destination; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_artifact_promoted_to_destination ON public.artifact_promotion_approval_request USING btree (artifact_id, destination_pipeline_id) WHERE ((status = 3) AND (status = 1));


--
-- Name: idx_unique_k8s_event_watcher_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_k8s_event_watcher_name ON public.k8s_event_watcher USING btree (name) WHERE (active = true);


--
-- Name: idx_unique_oci_registry_config; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_oci_registry_config ON public.oci_registry_config USING btree (docker_artifact_store_id, repository_type) WHERE (deleted = false);


--
-- Name: idx_unique_policy_name_policy_of_version; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_policy_name_policy_of_version ON public.global_policy USING btree (name, policy_of, version) WHERE (deleted = false);


--
-- Name: idx_unique_profile_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_profile_name ON public.infra_profile USING btree (name) WHERE (active = true);


--
-- Name: idx_unique_task_type_and_identifier_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_task_type_and_identifier_id ON public.devtron_resource_task_run USING btree (task_type, task_type_identifier);


--
-- Name: idx_unique_user_group_identifier; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_user_group_identifier ON public.user_group USING btree (identifier) WHERE (active = true);


--
-- Name: idx_unique_user_group_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_user_group_name ON public.user_group USING btree (name) WHERE (active = true);


--
-- Name: idx_unique_user_group_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_unique_user_group_user_id ON public.user_group_mapping USING btree (user_id, user_group_id);


--
-- Name: image_path_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX image_path_index ON public.image_path_reservation USING btree (image_path);


--
-- Name: image_scan_deploy_info_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX image_scan_deploy_info_unique ON public.image_scan_deploy_info USING btree (scan_object_meta_id, object_type);


--
-- Name: image_scan_execution_history_id_ix; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX image_scan_execution_history_id_ix ON public.image_scan_execution_result USING btree (image_scan_execution_history_id);


--
-- Name: pco_ci_artifact_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX pco_ci_artifact_id ON public.pipeline_config_override USING btree (ci_artifact_id);


--
-- Name: pco_pipeline_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX pco_pipeline_id_idx ON public.pipeline_config_override USING btree (pipeline_id);


--
-- Name: role_unique; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX role_unique ON public.roles USING btree (role, access_type, approver);


--
-- Name: unique_deployment_app_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX unique_deployment_app_name ON public.pipeline USING btree (deployment_app_name, environment_id, deleted) WHERE (deleted = false);


--
-- Name: unique_profile_platform_mapping; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX unique_profile_platform_mapping ON public.profile_platform_mapping USING btree (profile_id, platform) WHERE (active = true);


--
-- Name: unique_user_request_action; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX unique_user_request_action ON public.request_approval_user_data USING btree (user_id, approval_request_id, request_type);


--
-- Name: user_audit_user_id_ix; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_audit_user_id_ix ON public.user_audit USING btree (user_id);


--
-- Name: version_history_git_hash_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX version_history_git_hash_index ON public.installed_app_version_history USING btree (git_hash);


--
-- Name: api_token api_token_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_token
    ADD CONSTRAINT api_token_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: app_env_linkouts app_env_linkouts_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_env_linkouts
    ADD CONSTRAINT app_env_linkouts_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: app_env_linkouts app_env_linkouts_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_env_linkouts
    ADD CONSTRAINT app_env_linkouts_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: app_label app_label_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_label
    ADD CONSTRAINT app_label_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: app_level_metrics app_metrics_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_level_metrics
    ADD CONSTRAINT app_metrics_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: app_status app_status_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_status
    ADD CONSTRAINT app_status_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: app_status app_status_env_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_status
    ADD CONSTRAINT app_status_env_id_fkey FOREIGN KEY (env_id) REFERENCES public.environment(id);


--
-- Name: app_store_application_version app_store_application_version_app_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store_application_version
    ADD CONSTRAINT app_store_application_version_app_store_id_fkey FOREIGN KEY (app_store_id) REFERENCES public.app_store(id);


--
-- Name: app_store app_store_chart_repo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store
    ADD CONSTRAINT app_store_chart_repo_id_fkey FOREIGN KEY (chart_repo_id) REFERENCES public.chart_repo(id);


--
-- Name: app_store_version_values app_store_version_values_app_store_application_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store_version_values
    ADD CONSTRAINT app_store_version_values_app_store_application_version_id_fkey FOREIGN KEY (app_store_application_version_id) REFERENCES public.app_store_application_version(id);


--
-- Name: app app_team_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app
    ADD CONSTRAINT app_team_id_fkey FOREIGN KEY (team_id) REFERENCES public.team(id);


--
-- Name: app_workflow app_workflow_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_workflow
    ADD CONSTRAINT app_workflow_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: app_workflow_mapping app_workflow_mapping_app_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_workflow_mapping
    ADD CONSTRAINT app_workflow_mapping_app_workflow_id_fkey FOREIGN KEY (app_workflow_id) REFERENCES public.app_workflow(id);


--
-- Name: artifact_promotion_approval_request artifact_promotion_approval_request_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.artifact_promotion_approval_request
    ADD CONSTRAINT artifact_promotion_approval_request_artifact_id_fkey FOREIGN KEY (artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: artifact_promotion_approval_request artifact_promotion_approval_request_policy_evaluation_audit_id_; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.artifact_promotion_approval_request
    ADD CONSTRAINT artifact_promotion_approval_request_policy_evaluation_audit_id_ FOREIGN KEY (policy_evaluation_audit_id) REFERENCES public.resource_filter_evaluation_audit(id);


--
-- Name: artifact_promotion_approval_request artifact_promotion_approval_request_policy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.artifact_promotion_approval_request
    ADD CONSTRAINT artifact_promotion_approval_request_policy_id_fkey FOREIGN KEY (policy_id) REFERENCES public.global_policy(id);


--
-- Name: auto_remediation_trigger auto_remediation_trigger_k8s_event_watcher_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.auto_remediation_trigger
    ADD CONSTRAINT auto_remediation_trigger_k8s_event_watcher_id_fkey FOREIGN KEY (watcher_id) REFERENCES public.k8s_event_watcher(id);


--
-- Name: cd_workflow cd_workflow_ci_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow
    ADD CONSTRAINT cd_workflow_ci_artifact_id_fkey FOREIGN KEY (ci_artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: cd_workflow cd_workflow_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow
    ADD CONSTRAINT cd_workflow_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: cd_workflow_runner cd_workflow_runner_cd_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow_runner
    ADD CONSTRAINT cd_workflow_runner_cd_workflow_id_fkey FOREIGN KEY (cd_workflow_id) REFERENCES public.cd_workflow(id);


--
-- Name: cd_workflow_runner cd_workflow_runner_deployment_approval_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cd_workflow_runner
    ADD CONSTRAINT cd_workflow_runner_deployment_approval_request_id_fkey FOREIGN KEY (deployment_approval_request_id) REFERENCES public.deployment_approval_request(id);


--
-- Name: chart_category_mapping chart_category_mapping_app_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_category_mapping
    ADD CONSTRAINT chart_category_mapping_app_store_id_fkey FOREIGN KEY (app_store_id) REFERENCES public.app_store(id);


--
-- Name: chart_category_mapping chart_category_mapping_chart_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_category_mapping
    ADD CONSTRAINT chart_category_mapping_chart_category_id_fkey FOREIGN KEY (chart_category_id) REFERENCES public.chart_category(id);


--
-- Name: chart_env_config_override chart_env_config_override_chart_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_env_config_override
    ADD CONSTRAINT chart_env_config_override_chart_id_fkey FOREIGN KEY (chart_id) REFERENCES public.charts(id);


--
-- Name: chart_env_config_override chart_env_config_override_target_environment_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_env_config_override
    ADD CONSTRAINT chart_env_config_override_target_environment_fkey FOREIGN KEY (target_environment) REFERENCES public.environment(id);


--
-- Name: chart_group_deployment chart_group_deployment_chart_group_entry_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_deployment
    ADD CONSTRAINT chart_group_deployment_chart_group_entry_id_fkey FOREIGN KEY (chart_group_entry_id) REFERENCES public.chart_group_entry(id);


--
-- Name: chart_group_deployment chart_group_deployment_chart_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_deployment
    ADD CONSTRAINT chart_group_deployment_chart_group_id_fkey FOREIGN KEY (chart_group_id) REFERENCES public.chart_group(id);


--
-- Name: chart_group_deployment chart_group_deployment_installed_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_deployment
    ADD CONSTRAINT chart_group_deployment_installed_app_id_fkey FOREIGN KEY (installed_app_id) REFERENCES public.installed_apps(id);


--
-- Name: chart_group_entry chart_group_entry_app_store_application_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_entry
    ADD CONSTRAINT chart_group_entry_app_store_application_version_id_fkey FOREIGN KEY (app_store_application_version_id) REFERENCES public.app_store_application_version(id);


--
-- Name: chart_group_entry chart_group_entry_app_store_values_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_entry
    ADD CONSTRAINT chart_group_entry_app_store_values_version_id_fkey FOREIGN KEY (app_store_values_version_id) REFERENCES public.app_store_version_values(id);


--
-- Name: chart_group_entry chart_group_entry_chart_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chart_group_entry
    ADD CONSTRAINT chart_group_entry_chart_group_id_fkey FOREIGN KEY (chart_group_id) REFERENCES public.chart_group(id);


--
-- Name: charts charts_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charts
    ADD CONSTRAINT charts_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: charts charts_chart_repo_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.charts
    ADD CONSTRAINT charts_chart_repo_id_fkey FOREIGN KEY (chart_repo_id) REFERENCES public.chart_repo(id);


--
-- Name: ci_artifact ci_artifact_ci_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact
    ADD CONSTRAINT ci_artifact_ci_workflow_id_fkey FOREIGN KEY (ci_workflow_id) REFERENCES public.ci_workflow(id);


--
-- Name: ci_artifact ci_artifact_external_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact
    ADD CONSTRAINT ci_artifact_external_ci_pipeline_id_fkey FOREIGN KEY (external_ci_pipeline_id) REFERENCES public.external_ci_pipeline(id);


--
-- Name: ci_artifact ci_artifact_parent_ci_artifact_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact
    ADD CONSTRAINT ci_artifact_parent_ci_artifact_fkey FOREIGN KEY (parent_ci_artifact) REFERENCES public.ci_artifact(id);


--
-- Name: ci_artifact ci_artifact_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_artifact
    ADD CONSTRAINT ci_artifact_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: ci_env_mapping ci_env_mapping_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_env_mapping
    ADD CONSTRAINT ci_env_mapping_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: ci_env_mapping ci_env_mapping_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_env_mapping
    ADD CONSTRAINT ci_env_mapping_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: ci_pipeline ci_pipeline_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline
    ADD CONSTRAINT ci_pipeline_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: ci_pipeline ci_pipeline_ci_template_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline
    ADD CONSTRAINT ci_pipeline_ci_template_id_fkey FOREIGN KEY (ci_template_id) REFERENCES public.ci_template(id);


--
-- Name: ci_pipeline_history ci_pipeline_history_ci_pipeline_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_history
    ADD CONSTRAINT ci_pipeline_history_ci_pipeline_id_fk FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: ci_pipeline_material ci_pipeline_material_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_material
    ADD CONSTRAINT ci_pipeline_material_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: ci_pipeline_material ci_pipeline_material_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_material
    ADD CONSTRAINT ci_pipeline_material_git_material_id_fkey FOREIGN KEY (git_material_id) REFERENCES public.git_material(id);


--
-- Name: ci_pipeline_scripts ci_pipeline_scripts_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_pipeline_scripts
    ADD CONSTRAINT ci_pipeline_scripts_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: ci_template ci_template_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: ci_template ci_template_build_context_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_build_context_git_material_id_fkey FOREIGN KEY (build_context_git_material_id) REFERENCES public.git_material(id);


--
-- Name: ci_template ci_template_ci_build_config_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_ci_build_config_id_fkey FOREIGN KEY (ci_build_config_id) REFERENCES public.ci_build_config(id);


--
-- Name: ci_template ci_template_docker_registry_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_docker_registry_id_fkey FOREIGN KEY (docker_registry_id) REFERENCES public.docker_artifact_store(id);


--
-- Name: ci_template_history ci_template_git_material_history_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_history
    ADD CONSTRAINT ci_template_git_material_history_id_fkey FOREIGN KEY (git_material_id) REFERENCES public.git_material(id);


--
-- Name: ci_template ci_template_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template
    ADD CONSTRAINT ci_template_git_material_id_fkey FOREIGN KEY (git_material_id) REFERENCES public.git_material(id);


--
-- Name: ci_template_history ci_template_history_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_history
    ADD CONSTRAINT ci_template_history_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: ci_template_history ci_template_history_ci_template_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_history
    ADD CONSTRAINT ci_template_history_ci_template_id_fkey FOREIGN KEY (ci_template_id) REFERENCES public.ci_template(id);


--
-- Name: ci_template_history ci_template_history_docker_registry_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_history
    ADD CONSTRAINT ci_template_history_docker_registry_id_fkey FOREIGN KEY (docker_registry_id) REFERENCES public.docker_artifact_store(id);


--
-- Name: ci_template_override ci_template_override_build_context_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_override
    ADD CONSTRAINT ci_template_override_build_context_git_material_id_fkey FOREIGN KEY (build_context_git_material_id) REFERENCES public.git_material(id);


--
-- Name: ci_template_override ci_template_override_ci_build_config_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_override
    ADD CONSTRAINT ci_template_override_ci_build_config_id_fkey FOREIGN KEY (ci_build_config_id) REFERENCES public.ci_build_config(id);


--
-- Name: ci_template_override ci_template_override_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_override
    ADD CONSTRAINT ci_template_override_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: ci_template_override ci_template_override_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_template_override
    ADD CONSTRAINT ci_template_override_git_material_id_fkey FOREIGN KEY (git_material_id) REFERENCES public.git_material(id);


--
-- Name: ci_workflow ci_workflow_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_workflow
    ADD CONSTRAINT ci_workflow_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: cluster_accounts cluster_accounts_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_accounts
    ADD CONSTRAINT cluster_accounts_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: cluster_helm_config cluster_helm_config_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_helm_config
    ADD CONSTRAINT cluster_helm_config_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: cluster_installed_apps cluster_installed_apps_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_installed_apps
    ADD CONSTRAINT cluster_installed_apps_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: cluster_installed_apps cluster_installed_apps_installed_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster_installed_apps
    ADD CONSTRAINT cluster_installed_apps_installed_app_id_fkey FOREIGN KEY (installed_app_id) REFERENCES public.installed_apps(id);


--
-- Name: config_map_app_level config_map_app_level_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_app_level
    ADD CONSTRAINT config_map_app_level_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: config_map_env_level config_map_env_level_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_env_level
    ADD CONSTRAINT config_map_env_level_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: config_map_env_level config_map_env_level_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_env_level
    ADD CONSTRAINT config_map_env_level_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: config_map_history config_map_history_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_history
    ADD CONSTRAINT config_map_history_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: config_map_pipeline_level config_map_pipeline_level_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_pipeline_level
    ADD CONSTRAINT config_map_pipeline_level_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: config_map_pipeline_level config_map_pipeline_level_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_pipeline_level
    ADD CONSTRAINT config_map_pipeline_level_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: config_map_pipeline_level config_map_pipeline_level_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.config_map_pipeline_level
    ADD CONSTRAINT config_map_pipeline_level_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: cve_policy_control cve_policy_control_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cve_policy_control
    ADD CONSTRAINT cve_policy_control_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: cve_policy_control cve_policy_control_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cve_policy_control
    ADD CONSTRAINT cve_policy_control_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: cve_policy_control cve_policy_control_cve_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cve_policy_control
    ADD CONSTRAINT cve_policy_control_cve_store_id_fkey FOREIGN KEY (cve_store_id) REFERENCES public.cve_store(name);


--
-- Name: cve_policy_control cve_policy_control_env_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cve_policy_control
    ADD CONSTRAINT cve_policy_control_env_id_fkey FOREIGN KEY (env_id) REFERENCES public.environment(id);


--
-- Name: devtron_resource_object_dep_relations dep_mapping_component_schema_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_object_dep_relations
    ADD CONSTRAINT dep_mapping_component_schema_id_fk FOREIGN KEY (component_dt_res_schema_id) REFERENCES public.devtron_resource_schema(id);


--
-- Name: devtron_resource_object_dep_relations dep_mapping_dependency_schema_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.devtron_resource_object_dep_relations
    ADD CONSTRAINT dep_mapping_dependency_schema_id_fk FOREIGN KEY (dependency_dt_res_schema_id) REFERENCES public.devtron_resource_schema(id);


--
-- Name: deployment_app_migration_history deployment_app_migration_history_app_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_app_migration_history
    ADD CONSTRAINT deployment_app_migration_history_app_id_fk FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: deployment_app_migration_history deployment_app_migration_history_env_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_app_migration_history
    ADD CONSTRAINT deployment_app_migration_history_env_id_fk FOREIGN KEY (env_id) REFERENCES public.environment(id);


--
-- Name: deployment_approval_request deployment_approval_request_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_approval_request
    ADD CONSTRAINT deployment_approval_request_artifact_id_fkey FOREIGN KEY (ci_artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: deployment_approval_request deployment_approval_request_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_approval_request
    ADD CONSTRAINT deployment_approval_request_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: request_approval_user_data deployment_approval_user_data_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.request_approval_user_data
    ADD CONSTRAINT deployment_approval_user_data_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: deployment_config deployment_config_app_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_config
    ADD CONSTRAINT deployment_config_app_id_fk FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: deployment_config deployment_config_env_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_config
    ADD CONSTRAINT deployment_config_env_id_fk FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: deployment_group_app deployment_group_app_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group_app
    ADD CONSTRAINT deployment_group_app_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: deployment_group_app deployment_group_app_deployment_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group_app
    ADD CONSTRAINT deployment_group_app_deployment_group_id_fkey FOREIGN KEY (deployment_group_id) REFERENCES public.deployment_group(id);


--
-- Name: deployment_group deployment_group_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group
    ADD CONSTRAINT deployment_group_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: deployment_group deployment_group_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_group
    ADD CONSTRAINT deployment_group_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: deployment_status deployment_status_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_status
    ADD CONSTRAINT deployment_status_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: deployment_status deployment_status_env_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_status
    ADD CONSTRAINT deployment_status_env_id_fkey FOREIGN KEY (env_id) REFERENCES public.environment(id);


--
-- Name: deployment_template_history deployment_template_history_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_template_history
    ADD CONSTRAINT deployment_template_history_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: docker_registry_ips_config docker_registry_ips_config_docker_artifact_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.docker_registry_ips_config
    ADD CONSTRAINT docker_registry_ips_config_docker_artifact_store_id_fkey FOREIGN KEY (docker_artifact_store_id) REFERENCES public.docker_artifact_store(id);


--
-- Name: draft_version_comment draft_versions_relation_draft_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.draft_version_comment
    ADD CONSTRAINT draft_versions_relation_draft_version_id_fkey FOREIGN KEY (draft_version_id) REFERENCES public.draft_version(id);


--
-- Name: draft_version drafts_relation_draft_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.draft_version
    ADD CONSTRAINT drafts_relation_draft_id_fkey FOREIGN KEY (draft_id) REFERENCES public.draft(id);


--
-- Name: draft_version_comment drafts_relation_draft_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.draft_version_comment
    ADD CONSTRAINT drafts_relation_draft_id_fkey FOREIGN KEY (draft_id) REFERENCES public.draft(id);


--
-- Name: env_level_app_metrics env_level_env_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.env_level_app_metrics
    ADD CONSTRAINT env_level_env_id_fkey FOREIGN KEY (env_id) REFERENCES public.environment(id);


--
-- Name: env_level_app_metrics env_metrics_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.env_level_app_metrics
    ADD CONSTRAINT env_metrics_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: environment environment_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.environment
    ADD CONSTRAINT environment_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: ephemeral_container_actions ephemeral_container_actions_ephemeral_container_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ephemeral_container_actions
    ADD CONSTRAINT ephemeral_container_actions_ephemeral_container_id_fkey FOREIGN KEY (ephemeral_container_id) REFERENCES public.ephemeral_container(id);


--
-- Name: ephemeral_container_actions ephemeral_container_actions_performed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ephemeral_container_actions
    ADD CONSTRAINT ephemeral_container_actions_performed_by_fkey FOREIGN KEY (performed_by) REFERENCES public.users(id);


--
-- Name: ephemeral_container ephemeral_container_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ephemeral_container
    ADD CONSTRAINT ephemeral_container_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: external_ci_pipeline external_ci_pipeline_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_ci_pipeline
    ADD CONSTRAINT external_ci_pipeline_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: external_ci_pipeline external_ci_pipeline_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_ci_pipeline
    ADD CONSTRAINT external_ci_pipeline_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: external_link_identifier_mapping external_link_cluster_mapping_external_link_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_link_identifier_mapping
    ADD CONSTRAINT external_link_cluster_mapping_external_link_id_fkey FOREIGN KEY (external_link_id) REFERENCES public.external_link(id);


--
-- Name: external_link external_link_external_link_monitoring_tool_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.external_link
    ADD CONSTRAINT external_link_external_link_monitoring_tool_id_fkey FOREIGN KEY (external_link_monitoring_tool_id) REFERENCES public.external_link_monitoring_tool(id);


--
-- Name: app_store fk_app_store_docker_artifact_store; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.app_store
    ADD CONSTRAINT fk_app_store_docker_artifact_store FOREIGN KEY (docker_artifact_store_id) REFERENCES public.docker_artifact_store(id);


--
-- Name: cluster fk_cluster_remote_connection_config; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cluster
    ADD CONSTRAINT fk_cluster_remote_connection_config FOREIGN KEY (remote_connection_config_id) REFERENCES public.remote_connection_config(id);


--
-- Name: docker_artifact_store fk_docker_artifact_store_remote_connection_config; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.docker_artifact_store
    ADD CONSTRAINT fk_docker_artifact_store_remote_connection_config FOREIGN KEY (remote_connection_config_id) REFERENCES public.remote_connection_config(id);


--
-- Name: ci_workflow fk_image_path_reservation_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ci_workflow
    ADD CONSTRAINT fk_image_path_reservation_id FOREIGN KEY (image_path_reservation_id) REFERENCES public.image_path_reservation(id);


--
-- Name: infra_profile_configuration fk_infra_profile_configuration_profile_platform_mapping_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infra_profile_configuration
    ADD CONSTRAINT fk_infra_profile_configuration_profile_platform_mapping_id FOREIGN KEY (profile_platform_mapping_id) REFERENCES public.profile_platform_mapping(id);


--
-- Name: profile_platform_mapping fk_profile; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.profile_platform_mapping
    ADD CONSTRAINT fk_profile FOREIGN KEY (profile_id) REFERENCES public.infra_profile(id) ON DELETE CASCADE;


--
-- Name: generic_note_history generic_note_history_generic_note_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.generic_note_history
    ADD CONSTRAINT generic_note_history_generic_note_id_fkey FOREIGN KEY (note_id) REFERENCES public.generic_note(id) ON DELETE CASCADE;


--
-- Name: git_provider git_host_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_provider
    ADD CONSTRAINT git_host_id_fkey FOREIGN KEY (git_host_id) REFERENCES public.git_host(id);


--
-- Name: git_material git_material_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_material
    ADD CONSTRAINT git_material_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: git_material git_material_git_provider_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_material
    ADD CONSTRAINT git_material_git_provider_id_fkey FOREIGN KEY (git_provider_id) REFERENCES public.git_provider(id);


--
-- Name: git_material_history git_material_history_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_material_history
    ADD CONSTRAINT git_material_history_git_material_id_fkey FOREIGN KEY (git_material_id) REFERENCES public.git_material(id);


--
-- Name: git_web_hook git_web_hook_ci_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_web_hook
    ADD CONSTRAINT git_web_hook_ci_material_id_fkey FOREIGN KEY (ci_material_id) REFERENCES public.ci_pipeline_material(id);


--
-- Name: git_web_hook git_web_hook_git_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.git_web_hook
    ADD CONSTRAINT git_web_hook_git_material_id_fkey FOREIGN KEY (git_material_id) REFERENCES public.git_material(id);


--
-- Name: global_policy_history global_policy_history_global_policy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_policy_history
    ADD CONSTRAINT global_policy_history_global_policy_id_fkey FOREIGN KEY (global_policy_id) REFERENCES public.global_policy(id);


--
-- Name: global_policy_searchable_field global_policy_searchable_field_global_policy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_policy_searchable_field
    ADD CONSTRAINT global_policy_searchable_field_global_policy_id_fkey FOREIGN KEY (global_policy_id) REFERENCES public.global_policy(id);


--
-- Name: global_policy_searchable_field global_policy_searchable_field_searchable_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.global_policy_searchable_field
    ADD CONSTRAINT global_policy_searchable_field_searchable_key_id_fkey FOREIGN KEY (searchable_key_id) REFERENCES public.devtron_resource_searchable_key(id);


--
-- Name: image_comments image_comment_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_comments
    ADD CONSTRAINT image_comment_artifact_id_fkey FOREIGN KEY (artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: image_comments image_comment_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_comments
    ADD CONSTRAINT image_comment_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: image_scan_deploy_info image_scan_deploy_info_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_deploy_info
    ADD CONSTRAINT image_scan_deploy_info_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: image_scan_deploy_info image_scan_deploy_info_env_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_deploy_info
    ADD CONSTRAINT image_scan_deploy_info_env_id_fkey FOREIGN KEY (env_id) REFERENCES public.environment(id);


--
-- Name: resource_scan_execution_result image_scan_execution_history_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_scan_execution_result
    ADD CONSTRAINT image_scan_execution_history_id_fkey FOREIGN KEY (image_scan_execution_history_id) REFERENCES public.image_scan_execution_history(id);


--
-- Name: image_scan_execution_result image_scan_execution_result_cve_store_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_execution_result
    ADD CONSTRAINT image_scan_execution_result_cve_store_name_fkey FOREIGN KEY (cve_store_name) REFERENCES public.cve_store(name);


--
-- Name: image_scan_execution_result image_scan_execution_result_image_scan_execution_history_i_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_scan_execution_result
    ADD CONSTRAINT image_scan_execution_result_image_scan_execution_history_i_fkey FOREIGN KEY (image_scan_execution_history_id) REFERENCES public.image_scan_execution_history(id);


--
-- Name: release_tags image_tag_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.release_tags
    ADD CONSTRAINT image_tag_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: release_tags image_tag_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.release_tags
    ADD CONSTRAINT image_tag_artifact_id_fkey FOREIGN KEY (artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: image_tagging_audit image_tagging_audit_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_tagging_audit
    ADD CONSTRAINT image_tagging_audit_artifact_id_fkey FOREIGN KEY (artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: image_tagging_audit image_tagging_audit_updated_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_tagging_audit
    ADD CONSTRAINT image_tagging_audit_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: infra_profile_configuration infra_profile_configuration_profile_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infra_profile_configuration
    ADD CONSTRAINT infra_profile_configuration_profile_id_fkey FOREIGN KEY (profile_id) REFERENCES public.infra_profile(id);


--
-- Name: infrastructure_installation_versions infrastructure_installation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.infrastructure_installation_versions
    ADD CONSTRAINT infrastructure_installation_id_fkey FOREIGN KEY (infrastructure_installation_id) REFERENCES public.infrastructure_installation(id);


--
-- Name: installed_app_version_history installed_app_version_history_installed_app_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_app_version_history
    ADD CONSTRAINT installed_app_version_history_installed_app_version_id_fkey FOREIGN KEY (installed_app_version_id) REFERENCES public.installed_app_versions(id);


--
-- Name: installed_app_versions installed_app_versions_app_store_application_version_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_app_versions
    ADD CONSTRAINT installed_app_versions_app_store_application_version_id_fkey FOREIGN KEY (app_store_application_version_id) REFERENCES public.app_store_application_version(id);


--
-- Name: installed_app_versions installed_app_versions_installed_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_app_versions
    ADD CONSTRAINT installed_app_versions_installed_app_id_fkey FOREIGN KEY (installed_app_id) REFERENCES public.installed_apps(id);


--
-- Name: installed_apps installed_apps_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_apps
    ADD CONSTRAINT installed_apps_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: installed_apps installed_apps_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.installed_apps
    ADD CONSTRAINT installed_apps_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: intercepted_event_execution intercepted_events_auto_remediation_trigger_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercepted_event_execution
    ADD CONSTRAINT intercepted_events_auto_remediation_trigger_id_fkey FOREIGN KEY (trigger_id) REFERENCES public.auto_remediation_trigger(id);


--
-- Name: intercepted_event_execution intercepted_events_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercepted_event_execution
    ADD CONSTRAINT intercepted_events_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: module_resource_status module_resource_status_module_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.module_resource_status
    ADD CONSTRAINT module_resource_status_module_id_fkey FOREIGN KEY (module_id) REFERENCES public.module(id);


--
-- Name: notification_settings notification_setting_notification_rule_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_setting_notification_rule_id_fk FOREIGN KEY (notification_rule_id) REFERENCES public.notification_rule(id);


--
-- Name: notification_settings notification_settings_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: notification_templates notification_settings_event_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_settings_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES public.event(id);


--
-- Name: notification_settings notification_settings_event_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES public.event(id);


--
-- Name: notification_settings notification_settings_event_view_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_event_view_id_fkey FOREIGN KEY (view_id) REFERENCES public.notification_settings_view(id);


--
-- Name: notification_settings notification_settings_team_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notification_settings
    ADD CONSTRAINT notification_settings_team_id_fkey FOREIGN KEY (team_id) REFERENCES public.team(id);


--
-- Name: notifier_event_log notifier_event_log_event_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifier_event_log
    ADD CONSTRAINT notifier_event_log_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES public.event(id);


--
-- Name: oci_registry_config oci_registry_config_docker_artifact_store_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.oci_registry_config
    ADD CONSTRAINT oci_registry_config_docker_artifact_store_id_fkey FOREIGN KEY (docker_artifact_store_id) REFERENCES public.docker_artifact_store(id);


--
-- Name: pipeline pipeline_app_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline
    ADD CONSTRAINT pipeline_app_id_fkey FOREIGN KEY (app_id) REFERENCES public.app(id);


--
-- Name: pipeline pipeline_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline
    ADD CONSTRAINT pipeline_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: pipeline_config_override pipeline_config_override_cd_workflow_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_config_override
    ADD CONSTRAINT pipeline_config_override_cd_workflow_id_fkey FOREIGN KEY (cd_workflow_id) REFERENCES public.cd_workflow(id);


--
-- Name: pipeline_config_override pipeline_config_override_ci_artifact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_config_override
    ADD CONSTRAINT pipeline_config_override_ci_artifact_id_fkey FOREIGN KEY (ci_artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: pipeline_config_override pipeline_config_override_env_config_override_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_config_override
    ADD CONSTRAINT pipeline_config_override_env_config_override_id_fkey FOREIGN KEY (env_config_override_id) REFERENCES public.chart_env_config_override(id);


--
-- Name: pipeline_config_override pipeline_config_override_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_config_override
    ADD CONSTRAINT pipeline_config_override_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: pipeline pipeline_environment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline
    ADD CONSTRAINT pipeline_environment_id_fkey FOREIGN KEY (environment_id) REFERENCES public.environment(id);


--
-- Name: pipeline_stage pipeline_stage_cd_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage
    ADD CONSTRAINT pipeline_stage_cd_pipeline_id_fkey FOREIGN KEY (cd_pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: pipeline_stage pipeline_stage_ci_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage
    ADD CONSTRAINT pipeline_stage_ci_pipeline_id_fkey FOREIGN KEY (ci_pipeline_id) REFERENCES public.ci_pipeline(id);


--
-- Name: pipeline_stage_step_condition pipeline_stage_step_condition_condition_variable_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_condition
    ADD CONSTRAINT pipeline_stage_step_condition_condition_variable_id_fkey FOREIGN KEY (condition_variable_id) REFERENCES public.pipeline_stage_step_variable(id);


--
-- Name: pipeline_stage_step_condition pipeline_stage_step_condition_plugin_step_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_condition
    ADD CONSTRAINT pipeline_stage_step_condition_plugin_step_id_fkey FOREIGN KEY (pipeline_stage_step_id) REFERENCES public.pipeline_stage_step(id);


--
-- Name: pipeline_stage_step_variable pipeline_stage_step_file_reference_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_variable
    ADD CONSTRAINT pipeline_stage_step_file_reference_id_fkey FOREIGN KEY (file_reference_id) REFERENCES public.file_reference(id);


--
-- Name: pipeline_stage_step pipeline_stage_step_ref_plugin_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step
    ADD CONSTRAINT pipeline_stage_step_ref_plugin_id_fkey FOREIGN KEY (ref_plugin_id) REFERENCES public.plugin_metadata(id);


--
-- Name: pipeline_stage_step pipeline_stage_step_script_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step
    ADD CONSTRAINT pipeline_stage_step_script_id_fkey FOREIGN KEY (script_id) REFERENCES public.plugin_pipeline_script(id);


--
-- Name: pipeline_stage_step_variable pipeline_stage_step_value_constraint_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_variable
    ADD CONSTRAINT pipeline_stage_step_value_constraint_id_fkey FOREIGN KEY (value_constraint_id) REFERENCES public.value_constraint(id);


--
-- Name: pipeline_stage_step_variable pipeline_stage_step_variable_pipeline_stage_step_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_stage_step_variable
    ADD CONSTRAINT pipeline_stage_step_variable_pipeline_stage_step_id_fkey FOREIGN KEY (pipeline_stage_step_id) REFERENCES public.pipeline_stage_step(id);


--
-- Name: pipeline_status_timeline pipeline_status_timeline_cd_workflow_runner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline
    ADD CONSTRAINT pipeline_status_timeline_cd_workflow_runner_id_fkey FOREIGN KEY (cd_workflow_runner_id) REFERENCES public.cd_workflow_runner(id);


--
-- Name: pipeline_status_timeline pipeline_status_timeline_installed_app_version_history_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline
    ADD CONSTRAINT pipeline_status_timeline_installed_app_version_history_id_fkey FOREIGN KEY (installed_app_version_history_id) REFERENCES public.installed_app_version_history(id);


--
-- Name: pipeline_status_timeline_resources pipeline_status_timeline_resources_cd_workflow_runner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline_resources
    ADD CONSTRAINT pipeline_status_timeline_resources_cd_workflow_runner_id_fkey FOREIGN KEY (cd_workflow_runner_id) REFERENCES public.cd_workflow_runner(id);


--
-- Name: pipeline_status_timeline_resources pipeline_status_timeline_resources_installed_app_version_histor; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline_resources
    ADD CONSTRAINT pipeline_status_timeline_resources_installed_app_version_histor FOREIGN KEY (installed_app_version_history_id) REFERENCES public.installed_app_version_history(id);


--
-- Name: pipeline_status_timeline_sync_detail pipeline_status_timeline_sync_detail_cd_workflow_runner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline_sync_detail
    ADD CONSTRAINT pipeline_status_timeline_sync_detail_cd_workflow_runner_id_fkey FOREIGN KEY (cd_workflow_runner_id) REFERENCES public.cd_workflow_runner(id);


--
-- Name: pipeline_status_timeline_sync_detail pipeline_status_timeline_sync_detail_installed_app_version_hist; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_status_timeline_sync_detail
    ADD CONSTRAINT pipeline_status_timeline_sync_detail_installed_app_version_hist FOREIGN KEY (installed_app_version_history_id) REFERENCES public.installed_app_version_history(id);


--
-- Name: pipeline_strategy_history pipeline_strategy_history_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_strategy_history
    ADD CONSTRAINT pipeline_strategy_history_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: pipeline_strategy pipeline_strategy_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pipeline_strategy
    ADD CONSTRAINT pipeline_strategy_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: plugin_metadata plugin_metadata_plugin_parent_metadata_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_metadata
    ADD CONSTRAINT plugin_metadata_plugin_parent_metadata_id_fkey FOREIGN KEY (plugin_parent_metadata_id) REFERENCES public.plugin_parent_metadata(id);


--
-- Name: plugin_step_condition plugin_step_condition_condition_variable_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step_condition
    ADD CONSTRAINT plugin_step_condition_condition_variable_id_fkey FOREIGN KEY (condition_variable_id) REFERENCES public.plugin_step_variable(id);


--
-- Name: plugin_step_condition plugin_step_condition_plugin_step_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step_condition
    ADD CONSTRAINT plugin_step_condition_plugin_step_id_fkey FOREIGN KEY (plugin_step_id) REFERENCES public.plugin_step(id);


--
-- Name: plugin_step_variable plugin_step_file_reference_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step_variable
    ADD CONSTRAINT plugin_step_file_reference_id_fkey FOREIGN KEY (file_reference_id) REFERENCES public.file_reference(id);


--
-- Name: plugin_step plugin_step_plugin_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step
    ADD CONSTRAINT plugin_step_plugin_id_fkey FOREIGN KEY (plugin_id) REFERENCES public.plugin_metadata(id);


--
-- Name: plugin_step plugin_step_ref_plugin_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step
    ADD CONSTRAINT plugin_step_ref_plugin_id_fkey FOREIGN KEY (ref_plugin_id) REFERENCES public.plugin_metadata(id);


--
-- Name: plugin_step plugin_step_script_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step
    ADD CONSTRAINT plugin_step_script_id_fkey FOREIGN KEY (script_id) REFERENCES public.plugin_pipeline_script(id);


--
-- Name: plugin_step_variable plugin_step_variable_plugin_step_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_step_variable
    ADD CONSTRAINT plugin_step_variable_plugin_step_id_fkey FOREIGN KEY (plugin_step_id) REFERENCES public.plugin_step(id);


--
-- Name: plugin_tag_relation plugin_tag_relation_plugin_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_tag_relation
    ADD CONSTRAINT plugin_tag_relation_plugin_id_fkey FOREIGN KEY (plugin_id) REFERENCES public.plugin_metadata(id);


--
-- Name: plugin_tag_relation plugin_tag_relation_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plugin_tag_relation
    ADD CONSTRAINT plugin_tag_relation_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.plugin_tag(id);


--
-- Name: pre_post_cd_script_history pre_post_cd_script_history_pipeline_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pre_post_cd_script_history
    ADD CONSTRAINT pre_post_cd_script_history_pipeline_id_fkey FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: pre_post_ci_script_history pre_post_ci_script_history_ci_pipeline_scripts_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pre_post_ci_script_history
    ADD CONSTRAINT pre_post_ci_script_history_ci_pipeline_scripts_id_fkey FOREIGN KEY (ci_pipeline_scripts_id) REFERENCES public.ci_pipeline_scripts(id);


--
-- Name: registry_index_mapping registry_index_mapping_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.registry_index_mapping
    ADD CONSTRAINT registry_index_mapping_id_fkey FOREIGN KEY (scan_tool_id) REFERENCES public.scan_tool_metadata(id);


--
-- Name: resource_filter_audit resource_filter_audit_filter_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_filter_audit
    ADD CONSTRAINT resource_filter_audit_filter_id_fkey FOREIGN KEY (filter_id) REFERENCES public.resource_filter(id);


--
-- Name: resource_group_mapping resource_group_mapping_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_group_mapping
    ADD CONSTRAINT resource_group_mapping_fk FOREIGN KEY (resource_group_id) REFERENCES public.resource_group(id);


--
-- Name: resource_qualifier_mapping resource_qualifier_mapping_global_policy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_qualifier_mapping
    ADD CONSTRAINT resource_qualifier_mapping_global_policy_id_fkey FOREIGN KEY (global_policy_id) REFERENCES public.global_policy(id);


--
-- Name: role_group_role_mapping role_group_role_mapping_role_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_group_role_mapping
    ADD CONSTRAINT role_group_role_mapping_role_group_id_fkey FOREIGN KEY (role_group_id) REFERENCES public.role_group(id);


--
-- Name: role_group_role_mapping role_group_role_mapping_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_group_role_mapping
    ADD CONSTRAINT role_group_role_mapping_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: scan_step_condition_mapping scan_step_condition_mapping_condition_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_step_condition_mapping
    ADD CONSTRAINT scan_step_condition_mapping_condition_id_fkey FOREIGN KEY (scan_step_condition_id) REFERENCES public.scan_step_condition(id);


--
-- Name: scan_step_condition_mapping scan_step_condition_mapping_tool_step_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_step_condition_mapping
    ADD CONSTRAINT scan_step_condition_mapping_tool_step_id_fkey FOREIGN KEY (scan_tool_step_id) REFERENCES public.scan_tool_step(id);


--
-- Name: scan_tool_execution_history_mapping scan_tool_execution_history_mapping_result_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_execution_history_mapping
    ADD CONSTRAINT scan_tool_execution_history_mapping_result_id_fkey FOREIGN KEY (image_scan_execution_history_id) REFERENCES public.image_scan_execution_history(id);


--
-- Name: scan_tool_execution_history_mapping scan_tool_execution_history_mapping_scan_tool_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_execution_history_mapping
    ADD CONSTRAINT scan_tool_execution_history_mapping_scan_tool_id_fkey FOREIGN KEY (scan_tool_id) REFERENCES public.scan_tool_metadata(id);


--
-- Name: scan_tool_metadata scan_tool_metadata_plugin_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_metadata
    ADD CONSTRAINT scan_tool_metadata_plugin_id_fkey FOREIGN KEY (plugin_id) REFERENCES public.plugin_metadata(id);


--
-- Name: scan_tool_step scan_tool_step_scan_tool_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.scan_tool_step
    ADD CONSTRAINT scan_tool_step_scan_tool_id_fkey FOREIGN KEY (scan_tool_id) REFERENCES public.scan_tool_metadata(id);


--
-- Name: script_path_arg_port_mapping script_path_arg_port_mapping_script_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.script_path_arg_port_mapping
    ADD CONSTRAINT script_path_arg_port_mapping_script_id_fkey FOREIGN KEY (script_id) REFERENCES public.plugin_pipeline_script(id);


--
-- Name: ses_config ses_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ses_config
    ADD CONSTRAINT ses_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: slack_config slack_team_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.slack_config
    ADD CONSTRAINT slack_team_name_fkey FOREIGN KEY (team_id) REFERENCES public.team(id);


--
-- Name: timeout_window_resource_mappings timeout_window_configuration_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timeout_window_resource_mappings
    ADD CONSTRAINT timeout_window_configuration_id_fkey FOREIGN KEY (timeout_window_configuration_id) REFERENCES public.timeout_window_configuration(id) ON DELETE CASCADE;


--
-- Name: user_audit user_audit_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_audit
    ADD CONSTRAINT user_audit_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_deployment_request user_deployment_request_cd_workflow_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_deployment_request
    ADD CONSTRAINT user_deployment_request_cd_workflow_id_fk FOREIGN KEY (cd_workflow_id) REFERENCES public.cd_workflow(id);


--
-- Name: user_deployment_request user_deployment_request_ci_artifact_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_deployment_request
    ADD CONSTRAINT user_deployment_request_ci_artifact_id_fk FOREIGN KEY (ci_artifact_id) REFERENCES public.ci_artifact(id);


--
-- Name: user_deployment_request user_deployment_request_pipeline_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_deployment_request
    ADD CONSTRAINT user_deployment_request_pipeline_id_fk FOREIGN KEY (pipeline_id) REFERENCES public.pipeline(id);


--
-- Name: user_group_mapping user_group_mapping_user_group_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_group_mapping
    ADD CONSTRAINT user_group_mapping_user_group_id_fkey FOREIGN KEY (user_group_id) REFERENCES public.user_group(id);


--
-- Name: user_group_mapping user_group_mapping_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_group_mapping
    ADD CONSTRAINT user_group_mapping_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: user_roles user_roles_timeout_window_configuration_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_timeout_window_configuration_id_fkey FOREIGN KEY (timeout_window_configuration_id) REFERENCES public.timeout_window_configuration(id);


--
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_terminal_access_data user_terminal_access_data_cluster_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_terminal_access_data
    ADD CONSTRAINT user_terminal_access_data_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES public.cluster(id);


--
-- Name: user_terminal_access_data user_terminal_access_data_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_terminal_access_data
    ADD CONSTRAINT user_terminal_access_data_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: slack_config users_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.slack_config
    ADD CONSTRAINT users_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: users users_timeout_window_configuration_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_timeout_window_configuration_id_fkey FOREIGN KEY (timeout_window_configuration_id) REFERENCES public.timeout_window_configuration(id);


--
-- Name: webhook_event_data webhook_event_data_ghid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_event_data
    ADD CONSTRAINT webhook_event_data_ghid_fkey FOREIGN KEY (git_host_id) REFERENCES public.git_host(id);


--
-- PostgreSQL database dump complete
--

