import {Component, OnInit} from '@angular/core';
import {MatSnackBar} from "@angular/material";
import {Message} from "./terminal/terminal.component";
import { Observable } from "rxjs";
import * as format from 'date-fns/format';
import 'rxjs/add/operator/delay'
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  private ws;
  private ports = new Set();
  public repositories = [];
  public messages: Message[] = [];

  constructor(private snakeBar: MatSnackBar) {
  }

  ngOnInit(): void {
    this.addMessage([{Body: 'connecting to websocket...', Color: 'green'}]);
    this.ws = Observable.webSocket({
      url: 'ws://localhost:1111/ws',
      openObserver: {
        next: () => {
            this.addMessage([{Body: 'websocket connection is open!', Color: 'green'}]);
            this.addMessage(    [{Color: 'green', Body: 'Hi!, start by clicking "Fetch Data" button',}]);
        }
      },
      closeObserver: {
        next: () => {
            this.addMessage([{Body: 'websocket connection closed.', Color: 'red'}]);
        }
      }
    });


    this.ws.filter(e => !e.hasOwnProperty('ExposedPorts'))
        .bufferTime(1000)
        .filter(arr => arr.length)
      .subscribe(events => {
        this.addMessage(events);
      });

    this.ws.filter(e => e.hasOwnProperty('ExposedPorts'))
      .delay(200)
      .subscribe((msg) => {
        msg.ExposedPorts = Array.from(new Set(msg.ExposedPorts));
        msg.ExposedPorts.forEach((port) => this.ports.add(port));
        this.repositories.push(msg);
      });
  }

  addMessage(msgs: Message[]) {
    msgs = msgs.map((msg) => {
      if (!msg.Timestamp) {
        msg.Timestamp = Date.now();
      }
      msg.Timestamp = format(msg.Timestamp, 'H:mm:ss');
      return msg;
    });

    this.messages = [...this.messages, ...msgs];
  }

  fetchRepositories() {
    this.ws.next(JSON.stringify({type: 'GET_USED_PORTS'}));
  }

  refresh() {
    window.location.reload();
  }

  generatePort() {
    let port;
    while (!port || this.ports.has(port)) {
      port = Math.floor((Math.random() * 8999) + 1000);
    }

    this.snakeBar.open(`Available Port Found:  ${port}`, null, {
      duration: 7000
    });
  }
}
