package dev.palindrom615.eeaao;

import com.google.gradle.osdetector.OsDetector;
import org.gradle.api.Action;
import org.gradle.api.tasks.Exec;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;

/**
 * The task to execute eeaao-codegen-cli
 */
public class GenerateEeaaoTask extends Exec {
    private static final String EXECUTABLE_PATH_FMT = "dev/palindrom615/eeaao/eeaao-codegen-cli-%s-%s";
    private EeaaoCodegenOptions options = getProject().getExtensions().getByType(EeaaoCodegenOptions.class);
    private OsDetector osDetector = getProject().getExtensions().getByType(OsDetector.class);

    /**
     * Default constructor
     */
    public GenerateEeaaoTask() {
        setGroup("eeaao");
        setDescription("Generate Eeaao code");
        setExecutable(findEeaaoCodegenCliPath());
        setEeaaoCodegenArgs();
    }

    /**
     * Configures the Eeaao codegen task
     * @param configurableAction
     */
    public void eeaaoCodegen(Action<EeaaoCodegenOptions> configurableAction) {
        options = new EeaaoCodegenOptions();
        configurableAction.execute(options);
        setEeaaoCodegenArgs();
    }

    private void setEeaaoCodegenArgs() {
        args(
                "--specdir",
                options.getSpecDir(),
                "--codeletdir",
                options.getCodeletDir(),
                "--outdir",
                options.getOutDir()
        );
    }

    private String findEeaaoCodegenCliPath() {
        final String executablePath = String.format(
                EXECUTABLE_PATH_FMT,
                OsDetectorConverters.convertOs(osDetector.getOs()),
                OsDetectorConverters.convertArch(osDetector.getArch())
        );
        final URL resource = getClass().getClassLoader().getResource(executablePath);
        if (resource == null) {
            throw new RuntimeException("eeaao-codegen-cli executable not found");
        }
        final Path tmpPath;
        try {
            tmpPath = Files.createTempFile("eeaao-codegen-cli", null);
        } catch (IOException e) {
            throw new RuntimeException("Failed to create temporary eeaao-codegen-cli executable file", e);
        }
        try (final InputStream is = resource.openStream()) {
            Files.copy(is, tmpPath, StandardCopyOption.REPLACE_EXISTING);
        } catch (Exception e) {
            throw new RuntimeException("Failed to copy eeaao-codegen-cli executable to temp file", e);
        }
        final File tmpFile = tmpPath.toFile();
        if (!tmpFile.setExecutable(true)) {
            throw new RuntimeException("Failed to set eeaao-codegen-cli executable permission");
        }
        tmpFile.deleteOnExit();
        return tmpPath.toString();
    }
}
