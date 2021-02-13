load(":android_ndk.bzl", "ANDROID_NDK")
load("@bazel_tools//tools/build_defs/cc:action_names.bzl", "ACTION_NAMES")
load(
   "@bazel_tools//tools/cpp:cc_toolchain_config_lib.bzl",
   "feature",
   "flag_group",
   "flag_set",
   "tool_path",
)

def _impl_android_386(ctx):
    # TODO: make hermetic
    tool_paths = [
        tool_path(
            name = "gcc",
            path = ANDROID_NDK + "/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android30-clang",
        ),
        tool_path(
            name = "ld",
            path = ANDROID_NDK + "/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android-ld.gold",
        ),
        tool_path(
            name = "ar",
            path = ANDROID_NDK + "/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android-ar",
        ),
        tool_path(
            name = "cpp",
            path = ANDROID_NDK + "/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android30-clang++",
        ),
        tool_path(
            name = "gcov",
            path = "/bin/false",
        ),
        tool_path(
            name = "nm",
            path = "/bin/false",
        ),
        tool_path(
            name = "objdump",
            path = "/bin/false",
        ),
        tool_path(
            name = "strip",
            path = "/bin/false",
        ),
    ]
    features = [
        feature(
            name = "default_linker_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = [
                        ACTION_NAMES.cpp_link_executable,
                        ACTION_NAMES.cpp_link_dynamic_library,
                        ACTION_NAMES.cpp_link_nodeps_dynamic_library,
                    ],
                    flag_groups = ([
                        flag_group(
                            flags = [
                                "-lstdc++",
                            ],
                        ),
                    ]),
                ),
            ],
        ),
    ]
    return cc_common.create_cc_toolchain_config_info(
        ctx = ctx,
        cxx_builtin_include_directories = [
            ANDROID_NDK + "/toolchains/llvm/prebuilt/linux-x86_64/include",
            ANDROID_NDK + "/toolchains/llvm/prebuilt/linux-x86_64/sysroot/usr/include"
        ],
        features = features,
        toolchain_identifier = "local",
        host_system_name = "local",
        target_system_name = "android",
        target_cpu = "@platforms//cpu:x86_32",
        target_libc = "unknown",
        compiler = "android_386",
        abi_version = "unknown",
        abi_libc_version = "unknown",
        tool_paths = tool_paths,
    )

cc_toolchain_config_android_386 = rule(
    implementation = _impl_android_386,
    attrs = {},
    provides = [CcToolchainConfigInfo],
)
