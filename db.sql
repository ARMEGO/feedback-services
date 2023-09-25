-- Table: public.employee

-- DROP TABLE IF EXISTS public.employee;

CREATE TABLE IF NOT EXISTS public.employee
(
    id integer NOT NULL DEFAULT 'nextval('employee_id_seq'::regclass)',
    username character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT employee_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.employee
    OWNER to postgres;

-- Table: public.feedback

-- DROP TABLE IF EXISTS public.feedback;

CREATE TABLE IF NOT EXISTS public.feedback
(
    id integer NOT NULL DEFAULT 'nextval('feedback_id_seq'::regclass)',
    owner character varying COLLATE pg_catalog."default" NOT NULL,
    rating integer NOT NULL DEFAULT 0,
    comments text COLLATE pg_catalog."default" NOT NULL,
    assigned_by character varying COLLATE pg_catalog."default" NOT NULL DEFAULT 'admin'::character varying,
    assigned_on time with time zone NOT NULL,
    reviewed_by character varying COLLATE pg_catalog."default",
    CONSTRAINT feedback_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.feedback
    OWNER to postgres;