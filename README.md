```                                  
                                    
                               ,-.  
  ,----.                   ,--/ /|                                          
 /   /  \-.         .---.,--. :/ |  
|   :    :|        /. ./|:  : ' /   
|   | .\  .     .-'-. ' ||  '  /    
.   ; |:  |    /___/ \: |'  |  :                          absolutely rapid
'   .  \  | .-'.. '   ' .|  |   \   
 \   `.   |/___/ \:     ''  : |. \  
  `--'""| |.   \  ' .\   |  | ' \ \ 
    |   | | \   \   ' \ |'  : |--'  
    |   | :  \   \  |--" ;  |,'     
    `---'.|   \   \ |    '--'       
      `---`    '---"                                                    

```

# qwk ( Quick ) CLI Tool

Making everything you do - absolutely rapid. 

qwk (Quick) is a lightweight CLI automation and execution tool that simplifies running complex and repetitive commands. With an easy-to-use 'qwkfile', you can define your own custom commands, substitute arguments, and even chain multiple commands together. qwk offers a flexible and convenient way to manage your command-line tasks.

## Features

* Define custom commands using an intuitive and simple syntax in a qwkfile
* Run custom commands via the command line with qwk [command]
* Define global variables in the qwkfile or make use of positional arguments and messages in cli execution
* Chain multiple commands together for complex and repetitive tasks

## Installation

To install Qwk, you can instal the executable from the releases or build it from source. For the later, you need to have Go installed on your machine. Clone the Qwk repository and build the binary:

```bash
git clone https://github.com/drew-mcl/qwk.git
cd qwk
go build -o qwk
```

After building the binary, you can move it to your desired location and add it to your PATH.

## Usage

### qwkfile

Create a file named qwkfile in your working directory. This file will hold the custom commands you intend to use in your current directory. The structure of the file should follow this format:

```yaml
TODO
```

TODO: explain syntax

## Exmaple qwkfile

```yaml
TODO
```

### Running Commands

To run a custom command, simply type qwk followed by the command name and any arguments:

```bash
qwk greet Bigman wonderful
```

This command will output: 

```bash
Hello, Bigman! Today is a wonderful day.
```

If you don't provide the arguments, the command will use their default values:

```bash
qwk greet
```

This command will output:

```bash
Hello, Stranger! Today is a great day.
```
