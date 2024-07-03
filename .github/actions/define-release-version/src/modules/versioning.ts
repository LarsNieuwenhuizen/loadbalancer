import * as semver from 'semver';
import { git } from './git';

enum BumpType {
    major = 'major',
    minor = 'minor',
    patch = 'patch',
}

const methods = {
    newVersion: (commits: any, previousVersion: string): [BumpType, string] => {
        let result = BumpType.patch;
        let version = semver.inc(previousVersion, BumpType.patch) || "";
        commits.forEach((commit: string) => {
            const parts = commit.match(/(?<type>[a-z]*)(?<scope>\(.*\))?(?<breaking>!*): (?<message>.*)/);
            if (parts) {
                const type = parts.groups?.type;
                const breaking = parts.groups?.breaking === '!';

                if (breaking) {
                    result = BumpType.major;
                    version = semver.inc(previousVersion, BumpType.major) || "";
                    return;
                }
                if (type === "feat" || type === "feature") {
                    result = BumpType.minor;
                    version = semver.inc(previousVersion, BumpType.minor) || "";
                    return
                }
            }
        });
        return [result, `v${version}`];
    },
}

export async function getVersionInfo(): Promise<{previous: string, bumpType: string, newVersion: string}> {
    const latestTag: string = await git.executeGitCommand('describe --tags --abbrev=0');
    const previous: string = latestTag.trim();
    const commitsAfterLatestTag: any = await git.executeGitCommand("log " + previous + "..HEAD --oneline --format=%s");
    const commitArray = commitsAfterLatestTag.split('\n');

    const [bumpType, newVersion] = methods.newVersion(commitArray, previous);

    return {previous, bumpType, newVersion};
}
