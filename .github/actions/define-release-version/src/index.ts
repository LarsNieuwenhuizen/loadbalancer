import * as core from "@actions/core";
import * as versioning from "./modules/versioning";
import { env } from "process";

try {
    const info = await versioning.getVersionInfo();
    core.setOutput("previous_version", info.previous);
    core.setOutput("bump", info.bumpType);
    core.setOutput("new_version", info.newVersion);

    if (env.GITHUB_ACTIONS === "true") {
        core.summary
            .addHeading("Versioning Information :rocket:")
            .addTable([
                [{data: "Previous Version", header: true}, {data: "Bump Type", header: true}, {data: "New Version", header: true}],
                [info.previous, info.bumpType, info.newVersion]
            ])
            .write()
    }
} catch (error: any) {
    core.setFailed(error.message);
}
