# p2p-chat-app
How to run app ? <br/>
go build <br/>
Run in first terminal ./chat <br/>
Run in second terminal ./chat -d (use multi address printed on first terminal) <br/>

Concepts related to p2p chat : <br/>
1. we should append \n for every message as we are using it as delimiting character while reading it.<br/>
2. S.close() used for closing stream.<br/>
3. once a stream is opened , the object of stream should be used again and again untill it is closed. thats why we use for loop to stay on same stream for chat. <br/>
4. In order to speak to other peer , our peer should add the multiaddress of first peer to its peer table.<br/>
