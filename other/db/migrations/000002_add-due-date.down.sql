ALTER TABLE `assignment`
DROP `due_date`;

ALTER TABLE `submission`
DROP `is_late`;

-- Remove Seed

DELETE FROM `assignment`
WHERE id = 5 AND workspace_id = 1;

DELETE FROM `testcase`
WHERE id = 3 AND assignment_id = 5;
```


