#!/bin/bash

# Проверка наличия директории для хранения сборок, если её нет, создать
BUILD_DIR="./cmd/builts"
if [ ! -d "$BUILD_DIR" ]; then
    mkdir -p "$BUILD_DIR"
fi

# Переменные для имен файлов сборок
TMP_BUILD="tmp_streaming"
CURRENT_BUILD="streaming"
LEGACY_BUILD="legacy_streaming"

YELLOW_COLORED="\e[43;30m"
GREEN_COLORED="\e[42;30m"
RED_COLORED="\e[41;30m"
COLORED_END="\e[0m"

TAB="    - "
TAB_C="      "

echo -e "$YELLOW_COLORED Building of 'streaming' application started.$COLORED_END"
# Запуск сборки Go-приложения
go build -o "$BUILD_DIR/$TMP_BUILD" ./cmd/main.go
echo -e "$YELLOW_COLORED Building of 'streaming' application finished.$COLORED_END"

# Сборка прошла успешно, tmp_streaming создался
if [ -f "$BUILD_DIR/$TMP_BUILD" ]; then
    echo -e "$TAB $GREEN_COLORED Building of 'streaming' successfully complete, file $BUILD_DIR/$TMP_BUILD generated.$COLORED_END"

    echo -e "$YELLOW_COLORED Checking the current build $BUILD_DIR/$CURRENT_BUILD is EXISTS:$COLORED_END"
    # Проверяем, есть ли streaming актуальный билд
    if [ -f "$BUILD_DIR/$CURRENT_BUILD" ]; then
        echo -e "$TAB $GREEN_COLORED The current build $BUILD_DIR/$CURRENT_BUILD is EXISTS.$COLORED_END"
        echo -e "$TAB $YELLOW_COLORED Renaming the current build $BUILD_DIR/$CURRENT_BUILD to $BUILD_DIR/$LEGACY_BUILD.$COLORED_END"
        # Если да, переименовываем его в legacy
        mv "$BUILD_DIR/$CURRENT_BUILD" "$BUILD_DIR/$LEGACY_BUILD"

        echo -e "$YELLOW_COLORED Checking the current build $BUILD_DIR/$CURRENT_BUILD no longer EXISTS:$COLORED_END"
        if [ -f "$BUILD_DIR/$CURRENT_BUILD" ]; then
            echo -e "$TAB $RED_COLORED The current build $BUILD_DIR/$CURRENT_BUILD is EXISTS. Failed!$COLORED_END"
            exit 1
        else
            echo -e "$TAB $GREEN_COLORED The current build $BUILD_DIR/$CURRENT_BUILD is no longer EXISTS.$COLORED_END"
        fi

        echo -e "$TAB $YELLOW_COLORED Renaming the temporary build $BUILD_DIR/$TMP_BUILD to $BUILD_DIR/$CURRENT_BUILD.$COLORED_END"
        # Если steaming переименовался в legacy_streaming
        mv "$BUILD_DIR/$TMP_BUILD" "$BUILD_DIR/$CURRENT_BUILD" #то переименовываем tmp_streaming в streaming

        echo -e "$YELLOW_COLORED Checking the temporary build $BUILD_DIR/$TMP_BUILD no longer EXISTS:$COLORED_END"
        if [ -f "$BUILD_DIR/$TMP_BUILD" ]; then
            echo -e "$TAB $RED_COLORED The temporary build $BUILD_DIR/$TMP_BUILD is EXISTS. Failed!$COLORED_END"
            exit 1
        else
          echo -e "$TAB $GREEN_COLORED The temporary build $BUILD_DIR/$TMP_BUILD is no longer EXISTS.$COLORED_END"
        fi
    else
        echo -e "$TAB $YELLOW_COLORED The current build $BUILD_DIR/$CURRENT_BUILD is not EXISTS.$COLORED_END$GREEN_COLORED I guess that's first running, continue...$COLORED_END"
        echo -e "$TAB $YELLOW_COLORED Renaming the temporary build $BUILD_DIR/$TMP_BUILD to $BUILD_DIR/$CURRENT_BUILD.$COLORED_END"
        # Если steaming переименовался в legacy_streaming
        mv "$BUILD_DIR/$TMP_BUILD" "$BUILD_DIR/$CURRENT_BUILD" #то переименовываем tmp_streaming в streaming

        echo -e "$YELLOW_COLORED Checking the temporary build $BUILD_DIR/$TMP_BUILD no longer EXISTS:$COLORED_END"
        if [ -f "$BUILD_DIR/$TMP_BUILD" ]; then
            echo "$TAB $RED_COLORED The temporary build $BUILD_DIR/$TMP_BUILD is EXISTS. Failed!$COLORED_END"
            exit 1
        else
          echo -e "$TAB $GREEN_COLORED The temporary build $BUILD_DIR/$TMP_BUILD is no longer EXISTS.$COLORED_END"
        fi
    fi

    echo -e "$GREEN_COLORED Running the current build $BUILD_DIR/$CURRENT_BUILD:$COLORED_END"
    ./$BUILD_DIR/$CURRENT_BUILD
    if [ $? -eq 0 ]; then
        echo -e "$TAB $GREEN_COLORED The current build is successfully running.$COLORED_END"
    else
        echo -e "$TAB $RED_COLORED The current build running FAILED!$COLORED_END"
        # Running the previous versions
        echo -e "$YELLOW_COLORED Running the previous version:$COLORED_END"
        # Check the previous version is exists
        if [ -f "$BUILD_DIR/$LEGACY_BUILD" ]; then
            echo -e "$TAB $GREEN_COLORED The previous build is EXISTS.$COLORED_END"
            echo -e "$GREEN_COLORED Running the previous build $BUILD_DIR/$LEGACY_BUILD:$COLORED_END"
            ./$BUILD_DIR/$LEGACY_BUILD
            if [ $? -eq 0 ]; then
                echo -e "$TAB $GREEN_COLORED The previous build is successfully running.$COLORED_END"
            else
                echo -e "$TAB $RED_COLORED The current build running FAILED!$COLORED_END"
                exit 1
            fi
        fi
    fi
else
    echo -e "$TAB $RED_COLORED Building of 'streaming' FAILED! The temporary file was not generated.$COLORED_END"

    # Running the previous versions
    echo -e "$YELLOW_COLORED Running the previous versions:$COLORED_END"
    echo -e "$TAB_C $YELLOW_COLORED Checking the previous version is exists:$COLORED_END"
    # Check the previous version is exists
    if [ -f "$BUILD_DIR/$CURRENT_BUILD" ]; then
        echo -e "$TAB_C$TAB $GREEN_COLORED The previous version is EXISTS.$COLORED_END"
        echo -e "$TAB_C$TAB $YELLOW_COLORED Running the previous version $BUILD_DIR/$CURRENT_BUILD.$COLORED_END"
        ./$BUILD_DIR/$CURRENT_BUILD
    else
        echo -e "$TAB_C$TAB $RED_COLORED The previous version is NOT EXISTS!$COLORED_END"
        echo -e "$TAB_C $YELLOW_COLORED Checking the legacy version is exists:$COLORED_END"
        if [ -f "$BUILD_DIR/$LEGACY_BUILD" ]; then
            echo -e "$TAB_C$TAB $GREEN_COLORED The legacy version is EXISTS.$COLORED_END"
            echo -e "$TAB_C$TAB $YELLOW_COLORED Running the legacy version $BUILD_DIR/$LEGACY_BUILD.$COLORED_END"
            ./$BUILD_DIR/$LEGACY_BUILD
        else
            echo -e "$TAB_C$TAB $RED_COLORED The legacy version is NOT EXISTS!$COLORED_END"
            exit 1
        fi
    fi
fi