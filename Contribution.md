# Contribution Guidelines

## I have found a bug!
Please search the issue tracker if somebody else already filed a ticket for the problem. If yes, please amend your information to the ticket. Otherwise, create a new one with the following information.

 - the used version, tag, branch or commit of the APIKit
 - the content of the used OpenAPIv2 definition (YAML or JSON)
 - a description of the error (with logs or terminal output)
 - a description of the expected behaviour

We will process the bug ticket as soon as possible.

## I have found a security issue!
Please contact us via email (apikit-security@experienceone.com)

We will need the following information.

- the used version, tag, branch or commit of the APIKit
- a description of the security issue
- attack vectors
- a possible fix for the issue

## What do I need to know to help?
If you are looking to help to with a code contribution our project using Golang please first have a look at the technologies we used to build the APIKit.

- [APIKit documentation](/README.md)
- [OpenAPIv2 specification](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)
- [Jennifer code generator](https://github.com/dave/jennifer)

## How do I make a contribution?
Find an issue that you are interested in addressing or a feature that you would like to add.
Fork the repository associated with the issue to your local GitHub organization. This means that you will have a copy of the repository under `your-GitHub-username/repository-name`.

Clone the repository to your local machine using `git clone https://github.com/<your-GitHub-username/repository-name>`. Create a new branch for your fix using `git checkout -b <branch-name>`. Make the appropriate changes for the issue you are trying to address or the feature that you want to add.

Use `git add <paths-of-changed-files>` to add the file contents of the changed files to the "snapshot" git uses to manage the state of the project, also known as the index. Use `git commit -m "insert a short message of the changes made here"` to store the contents of the index with a descriptive message.

Push the changes to the remote repository using `git push origin <branch-name>`.
Submit a pull request to the upstream repository. Title the pull request with a short description of the changes made and the issue or bug number associated with your change. For example, you can title an issue like so "Added more log outputting to resolve #4352".
In the description of the pull request, explain the changes that you made, any issues you think exist with the pull request you made, and any questions you have for the maintainer. It's OK if your pull request is not perfect (no pull request is), the reviewer will be able to help you fix any problems and improve it!

Wait for the pull request to be reviewed by a maintainer. Make changes to the pull request if the reviewing maintainer recommends them. Celebrate your success after your pull request is merged!

## Where can I go for help?
If you need help, you can ask questions via E-Mail: info@experienceone.com.
