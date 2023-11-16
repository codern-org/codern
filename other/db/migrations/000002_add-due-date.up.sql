ALTER TABLE `assignment`
ADD `due_date` DATETIME NULL AFTER `level`;

ALTER TABLE `submission`
ADD `is_late` BOOLEAN DEFAULT 0 NOT NULL `submitted_at`;

-- Seeding

INSERT INTO `assignment`
    (id, workspace_id, name, description, detail_url, memory_limit, time_limit, level, created_at, updated_at, due_date)
VALUES (
    5,
    1,
    'Keeratikorn Saga I',
    'The first story of our Keeratikorn starts...',
    '/workspaces/1/assignments/5/detail/problem.md',
    1500,
    1000,
    'HARD',
    '2023-08-20 09:30:00',
    '2023-08-20 09:30:00',
    '2024-12-25 23:59:00'
);

INSERT INTO `testcase`
    (id, assignment_id, input_file_url, output_file_url)
VALUES (
    3,
    5,
    '/workspaces/1/assignments/5/testcase/1.in',
    '/workspaces/1/assignments/5/testcase/1.out'
);
