# services

A service must implement IService

- Start(ctx)
  This method is a non-blocking method.  It should spawn a go routine for its service.
  ctx is passed to the go routine　run(ctx)
- Stop
  -
- Status
  It should return status
  - nil running successfully (TBD)
  - error
- Reload
  It should gracefully reload configuration, and resume service.



Optionally the service can subscribe to signal service.
