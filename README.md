# nvim-go-client-examples
Examples of using the nvim [go client](https://github.com/neovim/go-client) and documentation for the client.

This file, the code, and the comments in the code are intended to help those who are interested in using Go to write
plugins for Neovim. Please do not take this document as absolute fact. I have done my best to document what I have
discovered by experimentation and reading the code, but it is very likely at least a few things I write here aren't 100%
accurate. I accept patches to the documentation and to the example code. My hope is to make this document as accurate
and as useful as possible.

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
The go-plugin code uses the default Go logger. In go-plugin/nvim/plugin/main.go there is a comment:
```
// Applications should use the default logger in the standard log package to
// write to Nvim's log.
```
So the first few lines in the example code's main() function do just this:
```
  // create a log to log to right away. It will help with debugging
  l, _ := os.Create("nvim-go-client-example.log")
  log.SetOutput(l)
```
This is pretty straightforward. Create a file to log to and set the output to log to that file. The
go-plugin code doesn't perform a lot of logging so I have found this mostly useful for debugging my
own code when using go-plugin. 

The rest of the code in main creates new functions, commands, and autocommands via Go. 

A slight tangent to editorialize:

> I'm not certain that there is great advantage in creating your own commands and autocommands via Go. It's more
> cumbersome than creating them via Vimscript. The real advantage of using go-plugin is to create Vim functions that are
> written in Go. You can then call them from your Vimscript created autocommands and commands. Your Go functions can
> perform computationally intensive tasks much faster than Vimscript and can be called as if they were native functions. 

Regardless of what I said above these examples show you how to create commands and autocommands using Go.

The code main calls plugin.Main().
```
  plugin.Main(func(p *plugin.Plugin) error {
    ...
  }
```
This method does a number of things. It creates the basic flags (-manifest and -location) using the flag package and
runs the passed in p function. The p function is expected to have code that creates handlers for commands, autocommands,
and functions. After the p function is run Main runs nvim.Serve(). This starts your plugin waiting to send and receive
to and from Neovim. The Serve() method blocks forever.

As you get more comfortable with using go-plugin you may not want to continue to use Main(). Main() seems more like a
convenience function that does the bare minimum a good host should do. The only problem with not using the Main() method
is that it uses several useful helper functions to write out the manifest. Unfortunately neither of these helper methods
are exported. So you either have to write new ones or copy and paste the ones in the source to your code.

#### Commands
The first thing the p function does is use p.HandleComamnd() to create a new command:
```
  p.HandleCommand(&plugin.CommandOptions{Name: "ExCmd", NArgs: "?", Bang: true, Eval: "[getcwd(),bufname()]"},
      func(args []string, bang bool, eval *cmdEvalExample) {
          log.Print("called command ExCmd")
          exCmd(p, args, bang, eval)
      })
```
p.HandleCommand() takes two arguments. The first is a pointer to a plugin.CommandOptions struct and the second
is a function that implements the functionality for the new command.

**A Note About the Function Return Value**
> The second argument to p.HandleCommand (or p.HandleAutocmd or p.HandleFunction, as we'll see later) requires a very
> specific return value. The return value **must** be one of the following:
> 1. no return value:  func m() {}
> 2. returns an error:  func m() (error) {}
> 3. returns a type and an error:  func m() ([]string, error) {}
>
> This is enforced by go-client. If you ever see an error in :messages that looks like "msgpack/rpc: handler return 
> must be (), (error) or (valueType, error)" it means the return signature you have for your function is wrong.

The plugin.CommandOptions record has many fields, all of which correspond directly to the arguments you can pass
to :command within Neovim (see :help :command for more info). You can look at the fields in the record by looking
in go-plugin/nvim/plugin/plugin.go. A quick listing of the fields in the struct are: Name, NArgs, Range, Count, Addr,
Bang, Register, Eval, Bar, Complete. This example doesn't cover every option, but should help you figure out how
to use those options should you need them.

For this example the plugin.CommandOptions struct assigns a name of 'ExCmd' to the command, specifies that the number
of arguments is 0 or 1 (that's what the "?" means), says that a bang ("!") is allowed, and has an eval section.

The name is the name a user will use from Neovim to call the command. So, in this case, we have defined a command
with the name of ExCmd. So, from Neovim, you can use :ExCmd. Give it a try. Fire up Neovim and run 
```
:ExCmd! hi
```
Nothing much will happen since the command only logs to the log file. So feel free to quit Neovim and look at the 
log file that was created in the same directory. It will probably look something like this:
```
Just entered a buffer
called command ExCmd
  Args to exCmd:
    arg1: hi
    bang: %!s(bool=true)
    cwd: /home/seth
    buffer: 
```
You can see that it logged your use of :ExCmd, logged the argument you gave it ("hi"), logged that you used a bang,
and also current directory and a buffer name (in this case the buffer name was empty because the buffer had no name).

If we look closer at the second argument to p.HandleCommand():
```
  func(args []string, bang bool, eval *cmdEvalExample) {
      log.Print("called command ExCmd")
      exCmd(p, args, bang, eval)
  })
```
we can see that the code is using an anonymous function to handle the arguments, log the use of 'ExCmd', and also 
call another function. The anonymous function is useful because it creates a closure that can be used to pass 
the plugin object to the exCmd function. But before we get to that let's talk about the arguments to the anonymous
function.

The first argument is 'args' and it is typed as a slice of strings. This is where the arguments to :ExCmd get placed.
So, when you ran :ExCmd earlier you passed "hi". This is an argument and it was passed in to this function as the
first element in the args slice.

The second argument is 'bang' and it is typed as a bool. It simply lets us know if an exclamation point was given when
:ExCmd is called. Above, when you typed :ExCmd!, you used a bang. So in that case bang would have been true.

The third argument is 'eval' and it is a pointer to a struct. That struct is defined in commands.go and looks like:
```
  type cmdEvalExample struct {
      Cwd     string `msgpack:",array"`
      Bufname string
  }
```
The fields in the struct match up with the expression given in the 'Eval' field for plugin.CommandOptions. Notice that
the 'Eval' field is vimscript surrounded by quotes. The vimscript is a list with two fields, each field being the result
of an expression. The first expression is getcwd() and the second expression is bufname(). Note the field tag in the
cmdEvalExample struct definition. This works hand-in-hand with what is in Eval.

How does go-plugin map the defined fields in plugin.CommandOptions to the arguments in the function? It appears to
assume that the function will take the argument in the order they are defined in the plugin.CommandOptions struct. 
The order of the definition in the struct is Name, NArgs, Range, Count, Addr, Bang, Register, Eval, Bar, Complete. So
if you define Name, NArgs, Range, Bang, Eval, and Bar the function signature will look like:
```
func(args []string, range string, bang bool, eval *struct, bar bool)
```
I haven't tried every possible combination but this seems to be the case. This is one part of the go-plugin code I 
haven't examined thoroughly yet.

#### Autocommands
The example code creates two new autocommands. It does this using p.HandleAutocmd like so:
```
  p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufEnter", Group: "ExmplNvGoClientGrp", Pattern: "*"},
      func() {
          log.Print("Just entered a buffer")
      })
  p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufAdd", Group: "ExmplNvGoClientGrp", Pattern: "*", Eval: "*"},
      func(eval *autocmdEvalExample) {
          log.Printf("buffer has cwd: %s", eval.Cwd)
      })
```
It's very similar to how commands are defined. Instead of calling p.HandleCommand it calls p.HandleAutocmd. Instead
of filling in a plugin.CommandOptions struct it fills in a plugin.AutocmdOptions struct (you can see the fields in
go-plugin/nvim/plugin/plugin.go). There are fewer fields for defining autocommands, but, like commands, they matchup
exactly with what you would use if you were defining the autocommand in vimscript. The fields available are:
Event, Group, Pattern, Nested, Eval. For more details on how autocommands are defined in vimscript see :help
autocommand.

My experimentation showed that if you define more than one autocommand in the same Event the manifest will be screwed
up. So two (or more) autocommands in event "BufEnter" seemed to cause a problem. However if one autocommand was in 
event "BufEnter,BufLeave" and the other was in "BufEnter" there were no problems. This brings up something that I 
don't show in the example code: the Event field can list more than one event. This maps exactly to how autocommands are
defined in vimscript, but may not be obvious from the examples I have provided.

The first autocommand created is very simple:
```
  p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufEnter", Group: "ExmplNvGoClientGrp", Pattern: "*"},
      func() {
          log.Print("Just entered a buffer")
      })
```
This will simply log that a buffer was entered into whenever that event is fired by Neovim. In fact, we have 
already seen this autocommand in action. If you go back and look at the log file excerpt in the previous section
you will see a line that says:
```
Just entered a buffer
```
That log entry was created by this autocommand.

The Group field is the same as the group you would use in a vimscript autocommand definition. The Pattern field is 
the same as the pattern you would pass to :autocommand in Neovim. i.e. it can be "\*.go" if you only want the autocommand
to work for Go files.

The second autocommand definition is:
```
  p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufAdd", Group: "ExmplNvGoClientGrp", Pattern: "*", Eval: "*"},
      func(eval *autocmdEvalExample) {
          log.Printf("buffer has cwd: %s", eval.Cwd)
      })
```
This definition includes an Eval section and also has an argument that is passed to the anonymous function. The Eval,
and how it works, is very similar to what was described above for the definition of the command. However, in this case,
the Eval is simply an asterisk. What does this mean? Let's look at the definition for the autocmdEvalExample struct:
```
type autocmdEvalExample struct {
	Cwd string `eval:"getcwd()"`
}
```
The field tag is now different. It now describes, directly, that the value of the given expression should be stored 
in the field. So, above, the value of the "getcwd()" vimscript expression should be stored in the Cwd field. You can
have multiple fields and each one can have different expressions. This is what the asterisk means in the Eval field
above. Here is an example of a struct with multiple fields, each field with its own expression:
```
type multiExprExample struct {
    Cwd string `eval:"getcwd()"`
    Name string `eval:"bufname()"`
    Id int `eval:"bufnr()"`
}
```
You can see the result from the second autocommand if you fire up Neovim and create a new buffer. You can then look
at the log file and see an entry that looks like:
```
Just entered a buffer
buffer has cwd: /home/seth
Just entered a buffer
```

#### Functions
The example code creates five new functions. It does this using p.HandleFunction like so:
```
  p.HandleFunction(&plugin.FunctionOptions{Name: "Upper"},
      func(args []string) (string, error) {
          log.Print("calling Upper")
          return upper(p, args[0]), nil
      })
  p.HandleFunction(&plugin.FunctionOptions{Name: "UpperCwd", Eval: "getcwd()"},
      func(args []string, dir string) (string, error) {
          log.Print("calling UpperCwd")
          return upper(p, dir), nil
      })
  p.HandleFunction(&plugin.FunctionOptions{Name: "ShowThings", Eval: "[getcwd(),argc()]"},
      func(args []string, eval *someArgs) ([]string, error) {
          log.Print("calling ShowThings")
          return returnArgs(p, eval)
      })
  p.HandleFunction(&plugin.FunctionOptions{Name: "GetVV"},
      func(args []string) ([]string, error) {
          log.Print("calling GetVV")
          return getvv(p, args[0])
      })
  p.HandleFunction(&plugin.FunctionOptions{Name: "ShowFirst"},
      func(args []string) (string, error) {
          log.Print("calling ShowFirst")
          return showfirst(p), nil
      })

```

At this point, if you have read the previous sections, it should be obvious what is going on. As it turns out, defining
functions is the simplest of the bunch. There are only two fields: Name and Eval (see the definition of
plugin.FunctionOptions in go-plugin/nvim/plugin/plugin.go). In this section I won't talk so much about Eval or the
definition of the functions. Instead there are a few things the functions do that may be interesting.

I continue to to use an anonymous function to wrap the code I want to use for the function. In this section we will
finally see a particularly good use for this: the ability of our code to use the plugin object.

The code above defines five new functions: Upper, UpperCwd, ShowThings, GetVV, and ShowFirst. The Go functions that each
of these functions call is in functions.go in the example code.

The "Upper" function simply takes a single string as argument, converts it to upper case and returns it. Feel free to
try it from within Neovim:

```
:echo Upper("this is now upper case")
```
Also, look at the log and notice the log entry:
```
calling Upper
```
Feel free to also run UpperCwd and ShowThings. Look at the code, see what each one does, and then run it.

The 'GetVV' function does something a little interesting. It uses the nvim API provided by go-plugin. The nvim API is
pretty large, but you can look at it using go doc:
```
go doc -all nvim
```
GetVV uses the VVar() method to retrieve whatever v: variable you wish to retrieve (for more information on v: variables
have a look at :help vim-variable). Perhaps the simplest one is v:oldfiles. So let's use GetVV to retrieve v:oldfiles.
Fire up Neovim and run:
```
:echo GetVV("oldfiles")
```
It should print out a fairly large Vim list. Each element in the list should be a path to a file that you have edited in
Neovim (somewhat) recently.

If we look at the code for GetVV:
```
func getvv(p *plugin.Plugin, name string) ([]string, error) {
	var result []string
	p.Nvim.VVar(name, &result)
	return result, nil
}
```
You can see that we use the plugin object. The plugin object has an Nvim object stashed in it and we use that to call
the VVar() method. The API for some of these methods is odd and may not be what you are expecting. It is definitely not
what I would have expected. The last argument to VVar() (per the VVar signature) is a pointer to interface{}. The
go-plugin restricts this to a pointer to either a slice or a map.

The last function defined is ShowFirst. This retrieves the first line of the current buffer and prints it out. This
demonstrates that your Go programs have access to the open buffers. So, first, try it out by bringing up Neovim and
running:
```
:echo ShowFirst()
```
You may want to type something in the empty buffer first. Or bring up a file with content on the first line. Otherwise
you'll just get a blank line when running ShowFirst().

The code used for ShowFirst is:
```
func showfirst(p *plugin.Plugin) string {
	br := nvim.NewBufferReader(p.Nvim, 0)
	r := bufio.NewReader(br)
	line, _ := r.ReadString('\n')
	return line
}
```
We once again use the plugin object. In this case we use it to create a new buffer reader object. See the go doc for the
nvim api for more information about the arguments to NewBufferReader(). Once we have the new buffer reader object it is
converted to a bufio.Reader and we can grab the first line of the buffer and return it. Pretty simple!
