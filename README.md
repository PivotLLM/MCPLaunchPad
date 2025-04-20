# MCPLaunchPad

This project servers as an example and a shell for creating an MCP server with a loosely coupled architecture that may need to be tweaked to fit your needs.

An interface and handler signature is defined in global/interface.go. The mcpserver package will accept any service that implements this interface.

The example1 and example2 packages implement the interface.

Package example1 demonstrates tools that communicates with an API. It includes some helper functions and assumes that the body returned by the API is in a format such as JSON that is suitable for direct use by the MCP client. It also assumes that the API is trusted not to return anything harmful.

Package example2 implements a simple tool that returns the time.

Note that the URL specified in main.go (api.example.com) does not exist, so every request provided by example1 will fail. Asking the LLM to get the time in 12 or 24 hour format (provided by example2) will work.

The mcpserver package calls the Register() function of each package to get the information required to register the tools the package provides.

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

## No Warranty (zilch, none, void, nil, null, "", {}, 0x00, 0b00000000, EOF)

THIS SOFTWARE IS PROVIDED “AS IS,” WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NON-INFRINGEMENT. IN NO EVENT SHALL THE COPYRIGHT HOLDERS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Made in Canada
