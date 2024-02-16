import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormsModule, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule,ReactiveFormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent implements OnInit {
  form: FormGroup;

  constructor(private fb: FormBuilder) {
      this.form = this.fb.group({
          email: ['', Validators.required],
          password: ['', Validators.required]
      });
  }

  ngOnInit() {
  }

  login() {
      console.log('Clicked the Login button');
  }
}