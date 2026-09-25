DO $$
DECLARE
	entity_table text;
BEGIN
	FOREACH entity_table IN ARRAY ARRAY[
		'Users',
		'FitnessProfile',
		'Target',
		'TrainingAvailability',
		'WorkoutPlan',
		'WorkoutSchedule',
		'Exercise',
		'WorkoutExercise',
		'WorkoutSet',
		'WorkoutSession',
		'ExercisePerformance',
		'BodyProgress',
		'RefreshToken'
	] LOOP
		IF to_regclass(format('%I', entity_table)) IS NOT NULL
		   AND EXISTS (
			   SELECT 1
			   FROM information_schema.columns
			   WHERE table_schema = current_schema()
				 AND table_name = entity_table
				 AND column_name = 'id'
				 AND udt_name = 'uuid'
				 AND column_default LIKE '%gen_random_uuid%'
		   ) THEN
			EXECUTE format(
				'ALTER TABLE %I ALTER COLUMN "id" DROP DEFAULT',
				entity_table
			);
		END IF;
	END LOOP;
END $$;
