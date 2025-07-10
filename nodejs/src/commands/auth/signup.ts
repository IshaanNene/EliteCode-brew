import { Command } from 'commander';
import inquirer from 'inquirer';
import chalk from 'chalk';
import { AuthService } from '../../services/auth';
import { validateEmail, validatePassword } from '../../utils/validators';

export const signupCommand = new Command()
  .name('signup')
  .description('Create a new Elitecode account')
  .action(async () => {
    try {
      const answers = await inquirer.prompt([
        {
          type: 'input',
          name: 'name',
          message: 'Name:',
          validate: (input) => input.length > 0 || 'Name is required'
        },
        {
          type: 'input',
          name: 'username',
          message: 'Username:',
          validate: (input) => input.length >= 3 || 'Username must be at least 3 characters'
        },
        {
          type: 'input',
          name: 'email',
          message: 'Email:',
          validate: validateEmail
        },
        {
          type: 'password',
          name: 'password',
          message: 'Password:',
          validate: validatePassword
        },
        {
          type: 'password',
          name: 'confirmPassword',
          message: 'Confirm Password:',
          validate: (input, answers) => {
            if (input !== answers.password) {
              return 'Passwords do not match';
            }
            return true;
          }
        }
      ]);

      const authService = new AuthService();
      const result = await authService.signup(answers);
      
      console.log(chalk.green('✅ Account created successfully!'));
      console.log(chalk.blue(`Welcome, ${result.user.name}!`));
      
    } catch (error) {
      console.error(chalk.red('❌ Signup failed:', error.message));
    }
  });