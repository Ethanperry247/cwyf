# Danger Dodgers Application

This is the front end component of the Danger Dodgers application. The software is built with react native -- currently it is only designed to work with Android, although ios support is a future goal of the project. 

TypeScript support may also be an ideal future addition to the Danger Dodgers front end. Currently, only javascript is supported.

## Project Structure

The project is currently organized in the following structure. Note that this structure overlooks some built-in react native files which should not be edited. 

```
DangerDodgersApp 
│
└───android
└───ios
└───src
│   └───Assets
│   └───Components
│   └───Helpers
│   └───Screens
│
│   App.js
│   config.js
│   index.js
```

- Android holds the generated gradle project which acts as the source code for an android app. Files within the android folder can be modified as needed. Examples of necessary modifications would be the inclusion of an additional user permission (which would go within the android manifest file), or any also custom/android-specific code which cannot be implemented with javascript on the react native side of the app. 
- Ios serves a similar purpose to android. Please note that this ios source code currently will not compile to a working app, or at the very least, it has not yet been tried. For a future ios side of this app, this would be the place to add user permissions as well as custom ios modules.
- Src contains JS react native code. This is where development should primarily occur. It is divided into four subdirectories for organizational purposes.
- Assets contains all static assets of the application. Examples would include large text files, images, and videos. 
- Components holds reusable react components.
- Helpers holds reusable helper functions and logic that is not a react component, meaning it doesn't contain any JSX.
- Screens holds any components which can stand as their own full screen view in the application. Examples would include a home page, the map page, or an account setting page. Any reusable logic which can be separated away from the screen components should be moved to the components folder.
- App.js contains the overall layout of the application. Consequentially, lots of logic makes its way into this file over time; however, it can tend to get out of control. Logic should be abstracted into the helper or components directories as separate functions and components.
- Index.js is the entrypoint for the application. Ideally, this file should not be edited, and any changes to the top level layout should occur within the App.js folder.

## How to Run this App

The best place to start with React Native development is getting the development environment running properly. A step-by-step guide for setting up this environment can be found on the react native website: [link](https://reactnative.dev/docs/environment-setup). 

It is unfortunately somewhat easy to make a mistake when setting up the development environment, so it is critical that the steps are following very carefully. At the conclusion of the environmental setup guide, android studio will have been set up and running. From this point, an actual android device or an emulator can be used to run the react native application. Please note that if developing with a local instance of the danger dodgers server running, this will affect the endpoint within the config.js file. If the emulator is being used, this endpoint should be set to 10.0.2.2. Otherwise if a physical device is being used, then the endpoint can be set to localhost, as long as ADB (android debug bridge) is running. 

A physical device is typically easiest to develop on, and is more realistic of the end result of any changes. If at all possible, it is recommended to use a physical device instead of an emulated one. 

## Requesting Resources from the Server

An authorized request context has been created within the App.js module and can be used to make authorized requests to the server/api. Usage looks like the following: 

```
const data = await context.useAuthorizedPost('/resource/', {
    objectProperty: {},
    someOtherProperty: 1
});
if (data.state === 'SUCCESS') {
    // Do something with the data
}
```

All other implementation details of sending the authorization token, or logging the user out if their user token timed out, are all handler by the methods in the authorization context. 