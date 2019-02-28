/*
 * Copyright 2017 NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.navercorp.pinpoint.plugin.rxjava;

import com.navercorp.pinpoint.bootstrap.config.ProfilerConfig;

/**
 * @author HyunGil Jeong
 */
public class RxJavaPluginConfig {

    private final boolean traceRxJava;

    public RxJavaPluginConfig(ProfilerConfig src) {
        this.traceRxJava = src.readBoolean("profiler.rxjava", false);
    }

    public boolean isTraceRxJava() {
        return traceRxJava;
    }
}
