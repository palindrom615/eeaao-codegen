package dev.palindrom615.eeaao;

/**
 * This class holds the options for the Eeaao codegen task.
 */
public class EeaaoCodegenOptions {
    private String codeletDir = "src/main/resources/codelet";
    private String outDir = "build/generated/sources/main/java";

    /**
     * @return the directory where the codelet files are located
     */
    public String getCodeletDir() {
        return codeletDir;
    }

    /**
     * @param codeletDir the directory where the codelet files are located
     */
    public void setCodeletDir(String codeletDir) {
        this.codeletDir = codeletDir;
    }

    /**
     * @return The directory where the generated code will be placed.
     */
    public String getOutDir() {
        return outDir;
    }

    /**
     * @param outDir The directory where the generated code will be placed.
     */
    public void setOutDir(String outDir) {
        this.outDir = outDir;
    }
}
