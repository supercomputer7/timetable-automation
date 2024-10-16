# A time table based automation

Modern UI for time table based automation, based on the wonderful [Semaphore UI](https://github.com/semaphoreui/semaphore), project.

I modified the original project, specifically for managing Linux IoT devices.

## What is the purpose of this project?

This project allows you to:
* Easily run local scripts as part of your daily automation tasks.
* Receive notifications about failed tasks.

## Key Concepts

1. **Projects** is a collection of related resources, configurations, and tasks. Each project allows you to organize and manage your automation efforts in one place, defining the scope of tasks such as deploying applications, running scripts, or orchestrating cloud resources. Projects help group resources, inventories, task templates, and environments for streamlined automation workflows.
2. **Task Templates** are reusable definitions of tasks that can be executed on demand or scheduled. A template specifies what actions should be performed, such as running Ansible playbooks, Terraform configurations, or other automation tasks. By using templates, you can standardize tasks and easily re-execute them with minimal effort, ensuring consistent results across different environments.
3. **Task** is a specific instance of a job or operation executed by Semaphore. It refers to running a predefined action (like an Ansible playbook or a script) using a task template. Tasks can be initiated manually or automatically through schedules and are tracked to give you detailed feedback on the execution, including success, failure, and logs.
4. **Schedules** allow you to automate task execution at specified times or intervals. This feature is useful for running periodic maintenance tasks, backups, or deployments without manual intervention. You can configure recurring schedules to ensure important automation tasks are performed regularly and on time.

## Getting Started

You can install this software by compiling it using the documentation which is provided
by the Semaphore UI project.

## Documentation (from the Semaphore UI project)

* [User Guide](https://docs.semaphoreui.com)
* [API Reference](https://semaphoreui.com/api-docs)

## License
This project is licensed under the [MIT license](https://github.com/supercomputer7/timetable-automation/LICENSE).
