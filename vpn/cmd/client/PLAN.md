# VPN Client Learning Plan

## Purpose

Build a Go command-line client that can establish a connection to a server, exchange data, report its state, and shut down cleanly.

This is a learning project. The early milestones are a TCP client and a tunnel protocol experiment, not a secure production VPN. Do not send Internet traffic through it or rely on it for privacy.

## Current starting point

The module is `example.com/vpnclient`.

The current command works as a flag-parsing exercise:

```bash
go run . connect --server 127.0.0.1:9000
```

`connectCmd` reads `--server` inside `RunE` with `cmd.Flags().GetString("server")`. Keep that pattern for command-specific flags. Positional arguments belong in `args`; named flags are obtained from `cmd`.

## Target command shape

The final learning CLI should have a small, deliberate surface:

```text
vpnclient connect --server host:port
vpnclient status
vpnclient disconnect
vpnclient version
```

At first, `status` and `disconnect` may report the state of the currently running process only. A later milestone can add a local state file or a background daemon. Do not invent persistence before a running connection exists.

## Boundaries

### In scope

- Cobra commands, flags, help text, errors, and exit codes.
- A TCP client connecting to a server you control.
- Context cancellation and Ctrl-C cleanup.
- Structured connection events and byte counters.
- TLS with a locally trusted development certificate.
- A deliberately simple framed protocol.
- UDP and TUN experiments after the TCP version works.

### Out of scope until explicitly studied

- A privacy or anonymity service.
- Production cryptography, key rotation, or user authentication.
- Routing arbitrary Internet traffic.
- DNS configuration or leak prevention.
- NAT traversal, anti-replay protection, path MTU discovery, and cross-platform TUN support.

## Phase 1: Make the CLI contract solid

### Goal

`connect` accepts a server address, validates it, and returns useful errors.

### Tasks

1. Replace Cobra's generated root-command text with an accurate name and short description.
2. Remove the generated `--toggle` flag. It is unrelated to the client.
3. Keep `connect --server` as a local flag on `connectCmd`.
4. Decide whether `--server` is required. For this client, it should be.
5. Validate the value as `host:port` with `net.SplitHostPort` before opening a connection.
6. Return errors from `RunE`; do not call `os.Exit` from a subcommand.
7. Add `--timeout`, using `time.Duration`, with a small default such as `5s`.
8. Add a `version` command that prints a fixed development version.

### Read before coding

- Cobra command lifecycle: `Use`, `Short`, `Args`, `RunE`, `Flags`, and `PersistentFlags`.
- `net.SplitHostPort` and `time.Duration`.
- `stdout`, `stderr`, and non-zero process exit codes.

### Done when

```bash
go run . connect --help
go run . connect --server invalid
go run . connect --server 127.0.0.1:9000 --timeout 3s
go run . version
```

produce understandable help or errors. Invalid input must fail before any network request.

## Phase 2: Establish one TCP connection

### Goal

Replace the placeholder `fmt.Println` with a real client-side TCP dial.

### Tasks

1. Create a context with the `--timeout` value.
2. Call `net.Dialer.DialContext(ctx, "tcp", server)`.
3. After connecting, log both endpoints:

   ```text
   local:  127.0.0.1:52341
   remote: 127.0.0.1:9000
   ```

4. Close the connection with `defer conn.Close()`.
5. Run against a local echo server first.
6. Distinguish a timeout, a refused connection, and an invalid address in the returned error message.

### Concepts to understand

`DialContext` creates the client-owned connected socket. The server has its own separate connected socket after calling `Accept`. The two sockets form one TCP connection identified by their local and remote address pair.

### Done when

The client connects to a local server, prints the two endpoint addresses, and closes without leaking a connection.

## Phase 3: Exchange bytes safely

### Goal

Understand that TCP transports a stream of bytes, not application messages.

### Tasks

1. Send one fixed greeting from the client.
2. Read an echo response.
3. Use `io.Reader` and `io.Writer` instead of tying code to `net.Conn` everywhere.
4. Add an interface only when it makes a boundary testable, for example:

   ```go
   type Connector interface {
       DialContext(context.Context, string, string) (net.Conn, error)
   }
   ```

5. Test the read/write logic with `net.Pipe` before relying on a real server.

### Important rule

Never assume one `Write` on one side equals one `Read` on the other. TCP has no message boundaries.

### Done when

The client can send and receive a known payload, including when the server sends that payload in multiple writes.

## Phase 4: Define packet framing

### Goal

Put application-level message boundaries on the TCP byte stream.

### First protocol

Use a small binary frame:

```text
4 bytes: payload length, unsigned big-endian
N bytes: payload
```

Set a maximum frame size before allocating memory. Start with 64 KiB.

### Tasks

1. Write `WriteFrame(w io.Writer, payload []byte) error`.
2. Write `ReadFrame(r io.Reader) ([]byte, error)`.
3. Return an error for zero-length policy violations, oversized lengths, and truncated frames.
4. Add tests for partial reads, two frames back-to-back, and maliciously large lengths.
5. Only after the framing tests pass, use the frame protocol from `connect`.

### Done when

Two frames sent consecutively are received as two distinct payloads regardless of how TCP chunks bytes.

## Phase 5: Concurrency and cancellation

### Goal

Make a live connection stop predictably when the user presses Ctrl-C or the server disconnects.

### Tasks

1. Use `signal.NotifyContext` for `SIGINT` and `SIGTERM`.
2. Run the send loop and receive loop in separate goroutines.
3. Give one owner responsibility for closing `net.Conn`.
4. Use a channel or `errgroup`-style coordination so the first fatal error ends both loops.
5. Ensure goroutines exit after cancellation; do not leave one blocked forever on a channel or socket read.
6. Print a single final result: disconnected normally, cancelled, or failed.

### State flow

```text
starting → dialing → connected → closing → disconnected
                     │
                     └→ failed
```

### Done when

Pressing Ctrl-C while connected closes the socket and terminates every client goroutine promptly.

## Phase 6: Add useful observability

### Goal

Make networking behavior visible without a debugger.

### Tasks

1. Use `log/slog` for connection lifecycle events.
2. Add `--verbose` for debug-level details.
3. Include connection ID, server address, elapsed time, bytes sent, bytes received, and error cause.
4. Never log payload contents, credentials, or private keys.
5. Keep the normal command output short; logs should explain failures when needed.

### Done when

You can answer: did it connect, to which address, for how long, how much data moved, and why did it close?

## Phase 7: Encrypt the transport with TLS

### Goal

Replace plain TCP with a TLS-protected connection to a development server.

### Tasks

1. Learn certificates, server names, certificate authorities, and the TLS handshake.
2. Generate a local development CA and a server certificate whose name matches the host you use.
3. Use `tls.Dialer` with a `tls.Config` that verifies the server certificate.
4. Do not set `InsecureSkipVerify: true` except as a short, explicit diagnostic experiment.
5. Report the negotiated TLS version and peer name only in verbose logs.

### Done when

The client rejects a server with an untrusted certificate and connects successfully when given the correct trusted CA.

## Phase 8: Decide what `status` and `disconnect` mean

### Goal

Choose a process model before adding commands that imply persistent state.

### Option A: Foreground client first

`connect` stays in the foreground. Ctrl-C disconnects it. `status` and `disconnect` are postponed or clearly state that no background session is managed.

### Option B: Local daemon later

`connect` asks a local background process to establish the tunnel. `status` and `disconnect` communicate with that process through a Unix socket.

Choose Option A first. It has fewer moving parts and makes connection ownership easy to understand.

## Phase 9: UDP transport experiment

### Goal

Learn how UDP differs from TCP without pretending UDP is a secure VPN.

### Tasks

1. Build a separate `udp-connect` command or a separate small program.
2. Use `net.UDPConn` and one datagram per framed payload.
3. Add read deadlines, a ping/pong keepalive, and a session identifier.
4. Simulate loss, duplication, reordering, and silence.
5. State clearly which delivery guarantees your protocol does and does not provide.

Do not reuse TCP assumptions: UDP preserves datagram boundaries but does not promise delivery, order, or uniqueness.

## Phase 10: Local-only TUN experiment

### Goal

Learn the boundary between Go and the operating system by carrying IP packets between two virtual interfaces.

### Preconditions

- The TCP/TLS/frame/cancellation milestones are working.
- You understand IPv4/IPv6 packet headers and MTU at a basic level.
- You are comfortable creating and removing routes manually.
- You use a disposable local environment or VM and elevated privileges only when required.

### Tasks

1. Create one TUN interface on each side.
2. Read an IP packet from the client TUN device.
3. Encapsulate it in the existing framed tunnel protocol.
4. Have the server write it to its TUN device; do the reverse for returned packets.
5. Route only a private test subnet through the interface.
6. Verify a destination with `route -n get <destination>` and inspect tables with `netstat -rn`.
7. Capture only your test traffic and observe encapsulated versus decapsulated packets.

Do not install a default route or change DNS while learning this stage.

## Suggested implementation order

Complete one small vertical slice at a time:

1. Clean CLI help and input validation.
2. TCP dialing and close behavior.
3. Send and receive one payload.
4. Framing with tests.
5. Cancellation with two concurrent loops.
6. Logging and counters.
7. TLS.
8. Foreground session behavior.
9. UDP experiment.
10. TUN and a private test route.

Do not start a later step until the previous one has a command you can run and an expected observable result.

## First task to do now

Stay in Phase 1. Replace the generated Cobra text, remove `--toggle`, make `--server` required, and validate it as `host:port`. Do not open a network connection yet.
