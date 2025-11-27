Objective
Your assignment is to implement a URL shortening service using Golang. You can use any framework you want (or no framework at all, it's up to you).

Brief
ShortLink is a URL shortening service where you enter a URL such as https://codesubmit.io/library/react and it returns a short URL such as http://your.domain/GeAi9K.

Tasks
Implement the assignment using:

Language: Golang
Two endpoints are required:
/encode: Encodes a URL to a shortened URL
/decode: Decodes a shortened URL to its original URL
Both endpoints should return JSON format
Simply ensure that a URL can be encoded into a short URL and that the short URL can be decoded back into the original URL.

Your application needs to be able to decode previously encoded URLs after a restart.

Provide detailed instructions on how to run your assignment in a separate markdown file.

Provide tests for both endpoints (and any other tests you may want to write).

You need to think through potential attack vectors on the application, and document them in the README.

You need to think through how your implementation may scale up, especially to solve the collision problem, and document your approach in the README.


Evaluation Criteria
Golang best practices
API implemented featuring a /encode and /decode endpoint
Completeness: Did you complete the features? Are all the tests running?
Correctness: Does the functionality act in sensible, thought-out ways?
Maintainability: Is it written in a clean, maintainable way?
Security: Have you identified potential issues and mitigated or documented them?
Scalability: What scalability issues do you foresee in your implementation and how do you plan to work around those issues?
Code Submit and Deployment
Please organize, design, test, and document your codebase as if it were going into a demo state, then push your code to GitHub. A public repo is okay.

After you have pushed your code, you may submit the assignment on the assignment page.

You can choose whatever free server to deploy your demo on.

All the best and happy coding.
Do not write any md files.