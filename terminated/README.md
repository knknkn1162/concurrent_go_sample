+ When parent send terminataion to child
    + parent should send `close(done)`
    + child should recv done channel and terminated

+ When parent send terminataion to child
  + child should send `close(done)`
  + parent should recv `<-done`.
