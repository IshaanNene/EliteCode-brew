#!/usr/bin/env node

import { Command } from 'commander';
import chalk from 'chalk';
import { authCommands } from './commands/auth';
import { problemCommands } from './commands/problems';
import { userCommands } from './commands/user';
import { systemCommands } from './commands/system';
import { githubCommands } from './commands/github';

const program = new Command();

program
  .name('elitecode')
  .description('CLI tool for competitive programming and coding challenges')
  .version('1.0.0');

// Add command groups
program.addCommand(authCommands);
program.addCommand(problemCommands);
program.addCommand(userCommands);
program.addCommand(systemCommands);
program.addCommand(githubCommands);

// Global error handler
program.exitOverride((err) => {
  console.error(chalk.red('Error:', err.message));
  process.exit(1);
});

// Parse command line arguments
program.parse(process.argv);

// Show help if no command provided
if (!process.argv.slice(2).length) {
  program.outputHelp();
}