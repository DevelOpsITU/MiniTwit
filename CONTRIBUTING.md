# Minitwit Contribution

## Which repository setup will we use?
We will be using multiple repositories in github. <br>
One named <b>MiniTwit</b> that contains the main application. This respository will hold all the source code for the web application. <br>
One named <b>MiniTwitDevelopment</b> that contains the docker containers for the web application, database and other applications that needs to be containerized. <br>
One named <b>ServerDeployment</b> that contains a make file and vagrant files to initialize the virtual box for the web hosting service and deploy the application.

## Which branching model will we use?
We will have a <b>main/master</b> branch, from which we create releases. The branch shall only hold a buildable and runnable application. All pushes to this branch will happen trough pull requests from <b>dev</b>.<br>
Secondly, the <b>dev</b> branch, will be the main developing branch. This branch should also only hold buildable and runnable code. <br>
From <b>dev</b>, other branches will be created per issue. These will be created under different subfolders, depending on the issue's topic. The general templates are seen listed below.
<ul>
    <li>features/#(ISSUEID)-(ISSUE_TITLE / RELEVANT_TITLE)</li>
    <li>bugs/#(ISSUEID)-(ISSUE_TITLE / RELEVANT_TITLE)
</ul>
All issue branches can only be merged into dev through pull-requests.

## Which distributed development workflow will we use?
We will use Github's Projects as an overview of all issues. From here, or directly from the seperate repositories, we will assign issues to different members. The member assigned is responsible for the issue, but may get help from the team. Most issue creation and closing will mostly happen weekly on Tuesdays after the lectures. However, issues may be given, taken, closed, opened and assigned during other days of the week by any member.

## How do we expect contributions to look like?
We expect members of the team to often commit, and be describtive in the commits. The commits' subject is limited to 50 chars, with few exceptions. If the commit has made changes to multiple files, due to forgetfulness of commiting in exitement of coding or other reasons, the member is recommended to put bulletpoints.

## Who is responsible for integrating/reviewing contributions?
We will strive to ensure that a given pull-request created by a given member is not reviewed by the same member.
