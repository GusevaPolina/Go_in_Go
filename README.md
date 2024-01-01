# Go in Go
The simple version of Go as a boarding game written in Go as a programming language. This is the experiment to test the possibilities of GPT-4 :neckbeard: :trollface: 

This is done in two paradigms: functional and OOP, respectively ->>

![Screenshot 2024-01-01 at 02 49 24](https://github.com/GusevaPolina/Go_in_Go/assets/107003133/55d83168-f8ff-4e4f-b814-d048d598c12d)
![image](https://github.com/GusevaPolina/Go_in_Go/assets/107003133/9845b1cc-c836-45f1-b970-83ab21deb8eb)

> $\color{red}\textnormal{Disclaimer}$: this code was created almost completely by ChatGPT 4. My main goal was to test if the subscription was worth it. I have never ever programmed in Go before. Moreover, GoLand by JetBrains was used as IDE instead of my preferred (and usually used) VS Code, party nights instead of calm office hours, and not importantly but MacOS instead of Windows :)

:hearts: :sunglasses: If you want to see the real stuff of Go in Go -> go to [the console version Go robots in Go made in 2009](https://github.com/skybrian/Gongo/tree/master) :gem:

#### Disclaimer mumba two: my take out of this experiment
This code is not perfect or even neat enough. It should not be. This is the code after three different versions of `please, clean` prompting. The only times I had to intervene were when creating the whole logic of cluster filling and minor debugging where ChatGPT got stuck. The most important part is that the minimal viable version was done in the first 5 hours. To get it into the nice app with more user-friendly features, it took 3 more hours. Thus, the simple but working app with a graphic interface and background music was composed with zero experience in the Go language in one day.

<br />

## How to 

### - install
To clone the repository: 
```bash
git clone https://github.com/GusevaPolina/Go_in_Go.git
```

To run (do not forget to download/install the compiler)
- one-pagers:
```bash
cd one_pagers
go run file_name.go    # for example, go run main_oop.go
```
- moduled versions
```bash
cd folder_name        # for example, cd modules_oop
go run .
```
### - play
The rules are simplified to the main one: take over as much place as possible! Only empty spaces could be conquered, no eye rule and etc.

To change the grid size, input a desired number in the input field.\
To reset, click the very left button (it's called `Delete all dotes` or `Go try again`).\
The timer resets after changing the grid size or resetting a game. The music will continue playing all the time.

The game is composed for two people playing: the first move is for blue dots and the second is for red ones.

Enjoy! :bowtie: :game_die:

