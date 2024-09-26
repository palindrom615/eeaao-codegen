package dev.palindrom615.eeaao;

/**
 * This class converts {@link com.google.gradle.osdetector.OsDetector}'s os and arch to that used in GOOS and GOARCH.
 * <p>
 * see respective specification on <a href="https://github.com/trustin/os-maven-plugin">os-maven-plugin</a> and <a href="https://golang.org/doc/install/source#environment">Golang</a>
 */
public final class OsDetectorConverters {
    private OsDetectorConverters() {
    }

    /**
     * Convert os to GOOS.
     * @param os given os-maven-plugin os
     * @return GOOS string
     */
    public static String convertOs(String os) {
        switch (os) {
            case "windows":
                return "windows";
            case "osx":
                return "darwin";
            case "linux":
                return "linux";
            case "freebsd":
                return "freebsd";
            case "openbsd":
                return "openbsd";
            case "netbsd":
                return "netbsd";
            case "sunos":
                return "solaris";
            default:
                throw new IllegalArgumentException("os not supported: " + os);
        }
    }

    /**
     * Convert arch to GOARCH.
     * @param arch given os-maven-plugin arch
     * @return GOARCH string
     */
    public static String convertArch(String arch) {
        switch (arch) {
            case "x86_32":
                return "386";
            case "x86_64":
                return "amd64";
            case "arm_32":
                return "arm";
            case "aarch_64":
                return "arm64";
            default:
                throw new IllegalArgumentException("architecture not supported: " + arch);
        }
    }
}
