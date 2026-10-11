--
-- PostgreSQL database dump
--


-- Dumped from database version 18.6
-- Dumped by pg_dump version 18.6


--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;




--
-- Name: body_progress; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.body_progress (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "BodyProgress_id_not_null" NOT NULL,
    "userId" uuid CONSTRAINT "BodyProgress_userId_not_null" NOT NULL,
    weight double precision,
    "recordedAt" timestamp without time zone,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: exercise; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.exercise (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "Exercise_id_not_null" NOT NULL,
    name character varying,
    "muscleGroup" character varying,
    instructions text,
    image character varying,
    video character varying,
    difficulty character varying,
    equipment character varying,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: exercise_performance; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.exercise_performance (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "ExercisePerformance_id_not_null" NOT NULL,
    "workoutSessionId" uuid CONSTRAINT "ExercisePerformance_workoutSessionId_not_null" NOT NULL,
    "workoutExerciseId" uuid CONSTRAINT "ExercisePerformance_workoutExerciseId_not_null" NOT NULL,
    "setNumber" integer,
    "actualReps" integer,
    "actualWeight" double precision,
    rpe double precision,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: fitness_profile; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.fitness_profile (
    "userId" uuid CONSTRAINT "FitnessProfile_userId_not_null" NOT NULL,
    gender character varying,
    age integer,
    height double precision,
    weight double precision,
    level character varying,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: refresh_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_token (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "RefreshToken_id_not_null" NOT NULL,
    "tokenHash" character varying NOT NULL,
    "tokenFamily" uuid NOT NULL,
    "userId" uuid CONSTRAINT "RefreshToken_userId_not_null" NOT NULL,
    "revokeAt" timestamp without time zone,
    "expiresAt" timestamp without time zone NOT NULL,
    "replacedById" uuid,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: target; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.target (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "Target_id_not_null" NOT NULL,
    "userId" uuid CONSTRAINT "Target_userId_not_null" NOT NULL,
    "targetWeight" double precision,
    goal character varying,
    calories double precision,
    protein double precision,
    carbs double precision,
    fat double precision,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: training_availability; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.training_availability (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "TrainingAvailability_id_not_null" NOT NULL,
    "userId" uuid CONSTRAINT "TrainingAvailability_userId_not_null" NOT NULL,
    "dayOfWeek" integer,
    "startTime" time without time zone,
    "endTime" time without time zone,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "User_id_not_null" NOT NULL,
    email character varying,
    "fullName" character varying,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone,
    password text
);


--
-- Name: workout_exercise; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workout_exercise (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "WorkoutExercise_id_not_null" NOT NULL,
    "workoutPlanId" uuid CONSTRAINT "WorkoutExercise_workoutPlanId_not_null" NOT NULL,
    "exerciseId" uuid CONSTRAINT "WorkoutExercise_exerciseId_not_null" NOT NULL,
    "orderIndex" integer,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: workout_plan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workout_plan (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "WorkoutPlan_id_not_null" NOT NULL,
    "userId" uuid CONSTRAINT "WorkoutPlan_userId_not_null" NOT NULL,
    "targetId" uuid CONSTRAINT "WorkoutPlan_targetId_not_null" NOT NULL,
    version integer DEFAULT 1,
    "durationPerSession" integer,
    "startDate" date,
    "endDate" date,
    status character varying,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: workout_schedule; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workout_schedule (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "WorkoutSchedule_id_not_null" NOT NULL,
    "workoutPlanId" uuid CONSTRAINT "WorkoutSchedule_workoutPlanId_not_null" NOT NULL,
    "dayOfWeek" integer,
    "startTime" time without time zone,
    duration integer,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: workout_session; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workout_session (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "WorkoutSession_id_not_null" NOT NULL,
    "userId" uuid CONSTRAINT "WorkoutSession_userId_not_null" NOT NULL,
    "workoutPlanId" uuid CONSTRAINT "WorkoutSession_workoutPlanId_not_null" NOT NULL,
    "scheduleId" uuid CONSTRAINT "WorkoutSession_scheduleId_not_null" NOT NULL,
    "startedAt" timestamp without time zone,
    "completedAt" timestamp without time zone,
    status character varying,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: workout_set; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.workout_set (
    id uuid DEFAULT gen_random_uuid() CONSTRAINT "WorkoutSet_id_not_null" NOT NULL,
    "workoutExerciseId" uuid CONSTRAINT "WorkoutSet_workoutExerciseId_not_null" NOT NULL,
    "setNumber" integer,
    "targetReps" integer,
    "targetWeight" double precision,
    "restTime" integer,
    "createdAt" timestamp without time zone,
    "updatedAt" timestamp without time zone
);


--
-- Name: body_progress BodyProgress_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.body_progress
    ADD CONSTRAINT "BodyProgress_pkey" PRIMARY KEY (id);


--
-- Name: exercise_performance ExercisePerformance_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exercise_performance
    ADD CONSTRAINT "ExercisePerformance_pkey" PRIMARY KEY (id);


--
-- Name: exercise Exercise_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exercise
    ADD CONSTRAINT "Exercise_pkey" PRIMARY KEY (id);


--
-- Name: fitness_profile FitnessProfile_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fitness_profile
    ADD CONSTRAINT "FitnessProfile_pkey" PRIMARY KEY ("userId");


--
-- Name: refresh_token RefreshToken_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT "RefreshToken_pkey" PRIMARY KEY (id);


--
-- Name: refresh_token RefreshToken_tokenHash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT "RefreshToken_tokenHash_key" UNIQUE ("tokenHash");


--
-- Name: target Target_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target
    ADD CONSTRAINT "Target_pkey" PRIMARY KEY (id);


--
-- Name: training_availability TrainingAvailability_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.training_availability
    ADD CONSTRAINT "TrainingAvailability_pkey" PRIMARY KEY (id);


--
-- Name: users User_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT "User_email_key" UNIQUE (email);


--
-- Name: users User_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (id);


--
-- Name: workout_exercise WorkoutExercise_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_exercise
    ADD CONSTRAINT "WorkoutExercise_pkey" PRIMARY KEY (id);


--
-- Name: workout_plan WorkoutPlan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_plan
    ADD CONSTRAINT "WorkoutPlan_pkey" PRIMARY KEY (id);


--
-- Name: workout_schedule WorkoutSchedule_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_schedule
    ADD CONSTRAINT "WorkoutSchedule_pkey" PRIMARY KEY (id);


--
-- Name: workout_session WorkoutSession_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_session
    ADD CONSTRAINT "WorkoutSession_pkey" PRIMARY KEY (id);


--
-- Name: workout_set WorkoutSet_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_set
    ADD CONSTRAINT "WorkoutSet_pkey" PRIMARY KEY (id);


--
-- Name: idx_refresh_token_family; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_token_family ON public.refresh_token USING btree ("tokenFamily");


--
-- Name: idx_refresh_token_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_token_user_id ON public.refresh_token USING btree ("userId");


--
-- Name: body_progress BodyProgress_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.body_progress
    ADD CONSTRAINT "BodyProgress_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: exercise_performance ExercisePerformance_workoutExerciseId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exercise_performance
    ADD CONSTRAINT "ExercisePerformance_workoutExerciseId_fkey" FOREIGN KEY ("workoutExerciseId") REFERENCES public.workout_exercise(id) DEFERRABLE;


--
-- Name: exercise_performance ExercisePerformance_workoutSessionId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.exercise_performance
    ADD CONSTRAINT "ExercisePerformance_workoutSessionId_fkey" FOREIGN KEY ("workoutSessionId") REFERENCES public.workout_session(id) DEFERRABLE;


--
-- Name: fitness_profile FitnessProfile_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fitness_profile
    ADD CONSTRAINT "FitnessProfile_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: refresh_token RefreshToken_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT "RefreshToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: target Target_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target
    ADD CONSTRAINT "Target_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: training_availability TrainingAvailability_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.training_availability
    ADD CONSTRAINT "TrainingAvailability_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: workout_exercise WorkoutExercise_exerciseId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_exercise
    ADD CONSTRAINT "WorkoutExercise_exerciseId_fkey" FOREIGN KEY ("exerciseId") REFERENCES public.exercise(id) DEFERRABLE;


--
-- Name: workout_exercise WorkoutExercise_workoutPlanId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_exercise
    ADD CONSTRAINT "WorkoutExercise_workoutPlanId_fkey" FOREIGN KEY ("workoutPlanId") REFERENCES public.workout_plan(id) DEFERRABLE;


--
-- Name: workout_plan WorkoutPlan_targetId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_plan
    ADD CONSTRAINT "WorkoutPlan_targetId_fkey" FOREIGN KEY ("targetId") REFERENCES public.target(id) DEFERRABLE;


--
-- Name: workout_plan WorkoutPlan_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_plan
    ADD CONSTRAINT "WorkoutPlan_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: workout_schedule WorkoutSchedule_workoutPlanId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_schedule
    ADD CONSTRAINT "WorkoutSchedule_workoutPlanId_fkey" FOREIGN KEY ("workoutPlanId") REFERENCES public.workout_plan(id) DEFERRABLE;


--
-- Name: workout_session WorkoutSession_scheduleId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_session
    ADD CONSTRAINT "WorkoutSession_scheduleId_fkey" FOREIGN KEY ("scheduleId") REFERENCES public.workout_schedule(id) DEFERRABLE;


--
-- Name: workout_session WorkoutSession_userId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_session
    ADD CONSTRAINT "WorkoutSession_userId_fkey" FOREIGN KEY ("userId") REFERENCES public.users(id) DEFERRABLE;


--
-- Name: workout_session WorkoutSession_workoutPlanId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_session
    ADD CONSTRAINT "WorkoutSession_workoutPlanId_fkey" FOREIGN KEY ("workoutPlanId") REFERENCES public.workout_plan(id) DEFERRABLE;


--
-- Name: workout_set WorkoutSet_workoutExerciseId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.workout_set
    ADD CONSTRAINT "WorkoutSet_workoutExerciseId_fkey" FOREIGN KEY ("workoutExerciseId") REFERENCES public.workout_exercise(id) DEFERRABLE;


--
-- PostgreSQL database dump complete
--


