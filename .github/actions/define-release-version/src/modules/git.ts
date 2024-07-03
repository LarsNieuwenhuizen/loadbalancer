import { exec } from 'child_process';
import { promisify } from "util";

export const git = {
    executeGitCommand: async function (command: string): Promise<string> {
      let gitCommand = `git ${command}`

      const execAsync = promisify(exec);
      return (await execAsync(gitCommand)).stdout;
    },
}
