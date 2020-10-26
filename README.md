# nvim-go-client-examples
Examples of using the nvim go client and documentation for the client

## Disclaimer(s)
* This is not intended to be a complete guide to all things [go-client](https://github.com/neovim/go-client). The api
provided by go-client is very large. Just have a look at the godoc for 'nvim' -- it's huge! 
* I show what I have learned. Particularly the stuff that struck me as non-obvious.
* If you want more examples or there is not an example of something you are interested in please feel free to do one of two things:
    1. Submit a bug report for this repository. If I can get to it, I will. 
    1. Figure it out yourself and submit a patch to this repository.

## Assumptions
* You are a Go programmer and have basic knowledge of Go modules and the Go programming language 
* You have a recent Go and Neovim installed
* You are basically familiar with Neovim project layout and installation of plugins

## Build & Install
You can clone the repository into the same place you place all your plugins. Lets assume this location is 
~/.config/nvim/pack/git-plugins/start. You will need to substitute wherever it is you install Vim plugins.

Clone the repository to where you install plugins:
```
cd ~/.config/nvim/pack/git-plugins/start
git clone https://github.com/WhoIsSethDaniel/nvim-go-client-examples
```
At this point if you start Neovim you will get some errors upon start. You will need to finish building and installing
the project before you can start Neovim without receiving errors.

Now you need to build the 'host'. In Neovim terms the 'host' is the program that will be talking to Neovim via
the msgpack RPC mechanism (see :help msgpack-rpc for more info). You don't actually need to know anything about the
msgpack protocol. If you want to learn more you can start [here](https://msgpack.org/index.html).

To build the project simply run make in the project directory
```
cd nvim-go-client-examples
make
```
This should build the 'nvim-go-client-example' program. This is the 'host' that will be talking to Neovim. By default
the host talks to Neovim over stdout/stdin. 

To verify the 'host' works as expected you should run:
```
./nvim-go-client-example -h
```
The output should be something like this:
```
Usage of ./nvim-go-client-example:
  -location .vim file
        Manifest is automatically written to .vim file
  -manifest host
        Write plugin manifest for host to stdout
```
You should copy the 'nvim-go-client-example' to somewhere in your path. Any time you rebuild the host you will need to
copy the new build to the same location. 

The host is now installed. When you start Neovim there should no longer be any errors.

## Code

### plugin/nvim-go-client.vim
This is the glue code that ties the Go 'host' to Neovim. This code registers the host, eventually starts the host, and
details what the host can do.

The code at the very bottom of the file that starts with 'call remote#host#RegisterPlugin' is copied and pasted from
running:
```
./nvim-go-client-example -manifest nvim_go_client_example
```
In other languages, such as Python, this step is somewhat automatic: you just run :UpdateRemotePlugins. For now the
Go client does not hook into this mechanism so you have to generate the manifest manually and paste it into your
Vim code.

The code just above that registers a Vim function that starts the host job. 
```
call remote#host#Register('nvim_go_client_example', 'x', function('s:Start_example_nvim_go_client'))
```
It does this by calling remote#host#Register(). The first argument to this function is the name of the host as given to
the -manifest argument above. e.g. nvim_go_client_example. The second argument, as best I can tell, can safely be
ignored. The third argument gives a function reference to a Vim function that starts the job. In this case that function
is s:Start_example_nvim_go_client().
```
function! s:panic(ch, data, ...) abort
    echom a:data
endfunction

function! s:Start_example_nvim_go_client(host) abort
    return jobstart(['nvim-go-client-example'], {
        \ 'rpc': v:true, 
        \ 'on_stderr': function('s:panic')
        \ })
endfunction
```
The function you register as the initiator of the host takes a single argument: the name of the host. In the example
in this code that argument is ignored. The key part of that function is the use of 'jobstart'. You can see :help
jobstart() for more information. The first argument to 'jobstart' is the name of the program we built earlier with
'make'. e.g. nvim-go-client-example. The 'rpc' argument is probably the most important since this tells Neovim that you 
want this program to use the msgpack RPC mechanism to communicate with Neovim. The 'on_error' argument is important
because it allows easier debugging of what is happening when things go wrong. Some errors get reported to stderr and all
this code does is make sure that that error gets printed to messages. See :help :messages for more information about
messages. Without the on_error section, and the function it calls, you will not see many of the errors that occur when
the msgpack encoding/decoding fails. Later on we can induce some errors and see what gets printed out.

### main.go

### functions.go

### autocmds.go

### commands.go
