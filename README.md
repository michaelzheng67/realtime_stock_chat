# Real time Stock Trading Chat App

`A full stack app that allows traders to look at stocks and chat in real time`


### Idea: 
In bloomberg terminal, you can get live data feeds of stocks as well as be able to chat with other traders. In this app, it provides users with an interface in which they can select the stock that they want to look at, as well as gives the user the ability to join a chat room associated with that stock where users can chat in real time about it.

![Screenshot 2023-11-07 at 10 04 19 PM](https://github.com/michaelzheng67/realtime_stock_chat/assets/82613778/0a129350-3f61-483e-94df-b30f7a196eba)


### Implementation:
When user opens up a tab on a given stock, a websocket connection is established with the stock data provider to retrieve values in real time. Then, when the user wants to open chat, another connection is opened with broadcaster running on localhost port, which lets users have group communication. Chat history is stored in a list while there is active users connected in order for new users to see the current chat history.

### Tech Stack:
Svelte frontend, Go backend


