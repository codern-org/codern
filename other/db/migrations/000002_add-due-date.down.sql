ALTER TABLE `assignment`
DROP `due_date`;

-- Remove Seed

DELETE FROM `testcase`
WHERE id = 3 AND assignment_id = 5;

DELETE FROM `assignment`
WHERE id = 5 AND workspace_id = 1;
