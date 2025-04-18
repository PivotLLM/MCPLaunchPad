# MCPLaunchPad

This project serves as an example and a shell for an MCP server. Tools can be added in mcpserver/tools.go

**CAUTION: This server is intended for local use and currently does not authenticate incoming requests because a standard interoperable mechanism for MCP clients to authenticate to MCP servers does not exist. It is hard-coded to listen on the localhost interface. Do not change this unless you are confident that you fully understand the risks and consequences.**

## Getting Started

### Prerequisites

- Go 1.24 or later

### Installation

```bash
# Clone the repository
git clone https://github.com/PivotLLM/MCPLaunchPad.git
cd MCPLaunchPad

# Build the project
go build -o mcplaunchpad
```

### Use

```bash
# Run the server with default settings
./mcplaunchpad

# For other options, see help
./mcplaunchpad -h
```

## Copyright and license

Copyright (c) 2025 by Tenebris Technologies Inc. This software is licensed under the MIT License. Please see LICENSE for details.

## No Warranty (nada, zilch, nil, null, "")

THIS SOFTWARE IS PROVIDED “AS IS,” WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NON-INFRINGEMENT. IN NO EVENT SHALL THE COPYRIGHT HOLDERS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
