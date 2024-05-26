
# Intro and porpoise

Serverless micro services with golang and aws cdk , and architecture overview. 

This is my weekend project. Microservices sounds cool, totally isolated services running to accomplish the greater goal, in theory scalable, event driven and clean. The big catch, that makes it impossible for medium and small projects is pricing. We’re in a time with compute is highly available  for everyone, free tiers are every time more generous. And yeah kubernettes is cool, kafka, redis, postgres. But each block you add sums cost, kafka can be around 60 usd, each service another 25, one database per service? sure! 60 each. And suddenly you’re spending 300 USD per month for an mvp nobody is using yet. 
And there my goal here, to create an affordable, scalable microservices system. What better than serverless? Aws for instance gives away 1M requests per month. If we combine it with SQS and the EventBus with dynamo db you got yourself a setup. And ohh lovely aws cdk, such a pleasure to structure all the infra with typescript, one microservice one class and provide all resources you need in a couple lines of code.

The source code is here if you want to take a look. My main focus was on developing a structure that I can reuse for other projects so there's no much of an usage yet but between `names`, `producer` and `hello` services we have a template for further developments. Next steps here would be to add authentication, some more tests and a real world usecase on top.


# Instructions





# TODO:

- Events routing
