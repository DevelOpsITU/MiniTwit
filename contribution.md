# Minitwit Contribution

## Which repository setup will we use?
We will be using multiple repositories in github. <br>
One named <b>MiniTwit</b> that contains the main application. This respository will hold all the source code for the web application. <br>
One named <b>miniTwitDevelopment</b> that contains the docker containers for the web application, database and other applications we find needed to be containerized. <br>
TODO: vagrant repo? 
## Which branching model will we use?
We will have a <b>main/master</b> branch which from we create releases. The branch shall only hold a buildable and runnable application. All pushes to this branch will happen trough pull requests from <b>dev</b>.<br>
Secondly, is the <b>dev</b> branch, that will be the main developing branch. This branch should also only hold buildable runnable code. <br>
From <b>dev</b> branches will be created per. issue, under different subfolders, depending on the issue's topic.
## Which distributed development workflow will we use?
## How do we expect contributions to look like?
## Who is responsible for integrating/reviewing contributions?


Which branching model will we use?
main
dev
feature-ISSUE-NAME

Which distributed development workflow will we use?
Work and talk with issues inmind
note conclusions in the issues

How do we expect contributions to look like?
commit often
limit subject to 50 chars
Maybe bullet points where it makes sense

To get code into Dev, must be a pull request with another developer as reviewer

Who is responsible for integrating/reviewing contributions?
Can git give it randomly?

Information to at least all of the above points should end up in a markdown document in your main repository (likely called CONTRIBUTE.md).