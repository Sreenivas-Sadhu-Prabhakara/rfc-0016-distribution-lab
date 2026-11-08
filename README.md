# RFC 16 — Memo distribution planner

**Covered RFC 16. Next RFC 17.** RFC 15's latest public guide, CLI source,
module and workflow were read before continuing. Its main branch was
`f92f34c7c863f86403036dcfb65ef85b8c8bc9ae` when checked.
Number 14 was recorded as Not Issued in the preceding session.

## Official document

Steve Crocker's **M.I.T.**, dated **27 August 1969**, adds MIT to the Network
Working Group memo distribution list. It supplies a recipient and physical
delivery details so MIT receives subsequent memos. This is administrative
coordination from the early ARPANET period.

The RFC Editor catalog was checked on **8 November 2026**. It lists stream
**Legacy** and publication/current status **Unknown**. RFC 16 **obsoletes RFC 10**
and is **updated by RFCs 24, 27 and 30**. Its own updates list and
obsoleted-by list are empty. These are the catalog's formal relationships,
distinct from a present-day recommendation to use historical contact details.

- [Official text](https://www.rfc-editor.org/rfc/rfc16.html)
- [RFC Editor record](https://www.rfc-editor.org/info/rfc16/)
- [Machine-readable catalog](https://www.rfc-editor.org/rfc/rfc16.json)

The modern architectural lesson is that participation changes need to reach
the people who distribute technical decisions. A correct memo can still fail
operationally if a newly responsible team is absent from the distribution list.

## Applications

1. **Banking API change notices, selected.** Compare subscriptions before an API
   change and require the gateway and banking teams to be present.
2. **Cloud incident roster review.** Identify added and removed notification
   destinations when ownership changes.
3. **RFC study registry.** Maintain a topic-based audience and preview the next
   distribution plan.

## Run

Requires Go 1.22 or later. From the repository root:

```sh
go run ./cmd/rosterlab < examples/roster.json
go run ./cmd/rosterlab -html < examples/roster.json > docs/index.html
go test -race ./...
go vet ./...
go build ./...
```

The example adds `banking@example.invalid`, removes `old-team@example.invalid`,
retains `gateway@example.invalid`, and plans two unique delivery addresses.
Two team IDs sharing the gateway mailbox produce one destination.

## Behavior and tradeoffs

This is a modern administrative tool inspired by the memo, with a synthetic
email roster replacing its historical postal distribution details.
The program only creates JSON or HTML reports. It has no delivery integration.

Input contains a selected `topic`, `previous` and `current` recipient arrays,
and optionally `required` team IDs. Each recipient has an `id`, `email`
and `topics`. Empty rosters are allowed. Required IDs must exist in the
current roster and subscribe to the selected topic.

The diff is by delivery endpoint. A contact correction produces a removed
old endpoint and an added new endpoint. Email domain case is normalized.
Local-part case is preserved, so `A@example.invalid` and
`a@example.invalid` remain distinct. Provider-specific aliases are not
inferred. This conservative choice can leave destinations that an individual
mail provider would treat as equivalent.

IDs and topic names contain 1..64 lowercase ASCII letters, digits or hyphens.
Each roster has at most 200 entries. Each entry has at most 32 unique topics.
Duplicate team IDs, malformed bare email addresses and duplicate topics fail
validation. Addresses accept up to 254 printable ASCII bytes with no spaces.
The parser performs syntax checks, with no DNS or mailbox existence checks.
The CLI rejects unknown JSON fields, trailing JSON and input over 1 MiB.
Failures exit nonzero before emitting a plan.

Output lists are sorted for review and diffing. The algorithm uses in-memory
sets sized to the bounded input. It favors repeatable previews over a persistent
subscription service. The hosted page displays one saved Go-generated plan.
The CLI remains the execution entry point.

## Validation

Race-enabled unit tests, vet and a local build passed using Go 1.23.12 on Linux.
Local compilation used `-buildvcs=false` for the scratch workspace.
The JSON and HTML entry points were executed against the sample.
Tests cover membership changes, alias deduplication, required-team omissions,
topic filtering, rejected rosters, contact corrections and CLI input bounds.
GitHub Actions runs tests, vet and build on repository pushes.
MIT licensed.
