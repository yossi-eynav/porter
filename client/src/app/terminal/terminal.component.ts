import { Component, Input, ViewChild} from '@angular/core';
import {throttle} from 'lodash';

export declare interface Message {
  Body: string,
  Color: string,
  Timestamp?: number,
}

@Component({
  selector: 'app-terminal',
  templateUrl: './terminal.component.html',
  styleUrls: ['./terminal.component.css']
})
export class TerminalComponent {
  private _logs: Message[] = [];
  @ViewChild('screen') screen;


  constructor() {
    this.scrollToBottom = throttle(this.scrollToBottom, 250);
  }

  scrollToBottom() {
    if (!this.screen) { return; }

    this.screen.nativeElement.scrollTo(0, this.screen.nativeElement.scrollHeight, { behavior: 'smooth'});
  }

  @Input()
  set logs(logs: Message[]) {
    this._logs = logs;
    setTimeout(() => {
      this.scrollToBottom();
    }, 0);
  }

  get logs() {
    return this._logs;
  }
}
