How long did this assignment take?
> According to Screen Time, I used GoLand for 2 hours, 39 minutes. That time does not include the time I spent
> working through a few tutorials and writing documentation. 

What was the hardest part?
> This was the first time I got to work with the **gorilla/mux** and **Gorm libraries**. 
> - Most of the Rest services I've written either used **net/http** directly or someone had already implemented a 
> lightweight framework that handled JWT authentication, method type validation, tracing token injection, etc. 
> - Most of the time, I've used **database/sql** when interacting with Postgres. 

Did you learn anything new?
> The **gorilla/mux** library cut down development time. The **Gorm** library eliminated the ORM code I was going to 
> need to write from scratch. I'm glad I got some exposure to these libraries. 

Is there anything you would have liked to implement but didn't have the time to?
> This version of the code is missing several things I wish I had time to correctly implement. 
> - Most of the microservices I've written have a wrapper that checks for a tracking token in the header and adds one 
> if it's not present. When using log aggregators like Splunk and DataDog, having a token you can use to track the path
> of a request through the system is very handy. 
> - Most of the service isn't sending back useful error codes. Specifically, there are places that should be returning
> 401, 403, 404, 405 or 409. Normally, services I've worked on use a custom error struct that has an error code so a 
> wrapper function can handle the response code.
> - I'm creating a Postgres connection per request instead using a connection pool. This can negatively impact
> performance and could run Postgres out of connections under heavy load.
> - I'm not using a context that times out when I call Postgres. Normally, I call *context.WithTimeout()* with the 
> Context provided by the *http.Request* to make sure callers can get a fail quickly when the DB is under heavy load. 
> - I'm not using a locking mutex when calling Postgres. Right now, competing requests could try to create a user with
> the same email address and get past the existing email check in *core.SignUp()*. 
> - I should be using a mutex in *core.UpdateUser()*. If a timestamp matching *updated_at* in Postgres would make it so 
> we could alert users when they're trying to update a user that has been updated since they originally called 
> *core.GetAllUsers()* but it would require changing the JSON sent to *core.UpdateUser()*. 
> - I didn't have time to write unit tests. 
> - There are several strings in the source code that I would have liked to change to properties in the config file. 

What are the security holes (if any) in your system? If there are any, how would you fix them?
> It's not a security hole per se, but you could DDOS attack the *core.SignUp()* and *core.SignIn()* methods. Adding a 
> rate limiter from **golang.org/x/time/rate** or using someone's custom rate limiter would have been a good idea, time 
> permitting. 

Do you feel that your skills were well tested?
> I think this assignment did a good job of covering what I knew about Rest services and showed how I normally write 
> them. 