import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';


import { AppComponent } from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatButtonModule, MatInputModule, MatListModule, MatSnackBarModule} from "@angular/material";
import { TerminalComponent } from './terminal/terminal.component';
import {CommonModule} from "@angular/common";

@NgModule({
  declarations: [
    AppComponent,
    TerminalComponent
  ],
  imports: [
    BrowserAnimationsModule,
    BrowserModule,
    MatListModule,
    MatButtonModule,
    MatInputModule,
    MatSnackBarModule,
    CommonModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
