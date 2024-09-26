package dev.palindrom615.eeaao;

import org.gradle.api.Plugin;
import org.gradle.api.Project;

/**
 * The plugin registers the default Eeaao codegen task.
 * <p>
 * For the plugin to support multiplatform, it applies {@link com.google.gradle.osdetector.OsDetectorPlugin} plugin.
 */
public class EeaaoCodegenPlugin implements Plugin<Project> {
    private static final String OSDETECTOR_PLUGIN = "com.google.osdetector";

    @Override
    public void apply(Project project) {
        project.getPluginManager().apply(OSDETECTOR_PLUGIN);

        project.getExtensions().create("eeaaoCodegen", EeaaoCodegenOptions.class);
        project.getTasks().register("eeaaoCodegen", GenerateEeaaoTask.class);
    }
}
