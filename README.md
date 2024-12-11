# go-bootcamp
Geekbang Go Basic Bootcamp

## v0.0.3

![kube](./doc/images/v0.0.3/kube.png)

- [X] Rate limiter middleware
- [ ] Unit tests
- [ ] Sync -> Async
  - [ ] Average response time
    - [ ] absolute value
    - [ ] trend (%x increase in 1 second)
  - [ ] Error rate (X% during some time)
  - [ ] Stop async mode
    - [ ] After N minutes
    - [ ] keep a percentage of sync request (using random number), if the sync request is ok (response time, no error), then gradually increase the sync percentage
  - [ ] Preempt send : Select for Update, only one instance can succeed to update the update time, then it continues to send the message, other instances will not query this record because it is now too recent.
