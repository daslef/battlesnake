import { useState, useEffect } from "react";

// import { playbackError, playbackState } from "../stores/playback";

import { useSettingsStore } from "../lib/stores/settings";

// import Gameboard from "$lib/components/Gameboard.svelte";
// import PlaybackControls from "$lib/components/PlaybackControls.svelte";
// import Scoreboard from "$lib/components/Scoreboard.svelte";


export default function Board() {
    const settingsStore = useSettingsStore()

    const [settingError, setSettingError] = useState(false)
    const [playbackError, setPlaybackError] = useState(false)

    useEffect(() => {
        const url = new URL(window.location.href);
        settings = loadSettingsWithURLOverrides(url);

        if (settings.game.length && settings.engine.length) {
            setSettingError(false);
            playbackState.load(settings);
        }

    }, [])

    return (
        <div class="w-full max-w-screen-xl md:aspect-video mx-auto">
            <div
                class="h-full w-full flex flex-col items-center justify-center"
            >
                {settingError ? (
                    <p class="p-4 font-bold text-lg text-center">
                        To display a game you need to specify the ID in the URL.
                        <p class="italic">
                            {`<url>?game=<GAME_ID>`}
                        </p>

                    </p>
                ) : playbackError ? (
                    <p class="p-4 font-bold text-lg text-center text-red-500">
                        {playbackError}
                    </p>
                ) : playbackState ? (
                    <div class="w-full h-full flex flex-col md:flex-row">
                        <div class="flex flex-col grow">
                            {settings.title &&
                                <h1 class="text-center font-bold pt-2 text-lg">{settings.title}</h1>
                            }
                            <Gameboard showCoordinates={settings.showCoords} />
                            {settings.showControls && (
                                <div class="flex justify-evenly text-xl py-2 px-6">
                                    <PlaybackControls />
                                </div>
                            )}
                        </div>
                        {settings.showScoreboard && (
                            <div class="basis-full md:basis-[45%] order-first p-2 md:order-last">
                                <Scoreboard />
                            </div>
                        )}
                    </div>
                ) :
                    <p class="p-4 text-lg text-center">Loading game...</p>
                }
            </div>
        </div>
    )
}