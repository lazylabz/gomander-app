import {
	createRecordingTerminals,
	type RecordingTerminals,
} from "@/commandOutput/adapters/recording.ts";
import {
	resetCommandOutputPipeline,
	resetTerminalFactory,
	setTerminalFactory,
} from "@/commandOutput/commandOutput.ts";

export const installRecordingTerminals = (): RecordingTerminals => {
	resetCommandOutputPipeline();
	const recording = createRecordingTerminals();
	setTerminalFactory(recording.create);
	return recording;
};

export const resetTerminals = (): void => {
	resetCommandOutputPipeline();
	resetTerminalFactory();
};
