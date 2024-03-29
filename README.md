# hpcadmin-server-go

## Comparison with Coldfront

Features we want:

- Sync with AD including delegated access for PIRG management
- Storage allocations tied to quotas
- Allow people to request access to pirgs
- and more

### Coldfront:

Pros:

- already has pirg/admins/members
- dump slurm associations using CLI
- supported by buffalo, though development is slow

Cons:

- no web api
  - automation is done via CLI load/dump commands
- no AD functionality
  - would need to write a custom submodule to dump associations
- slurm
  - assumes you're consuming association data generated by coldfront
- not easily extensible
  - adding features would either be a custom "application" (django term) or fork
    (bad idea)

### HPCAdmin

Pros:

- flexibility to build whatever we need

Cons:

- time
