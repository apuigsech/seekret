/*

Package seekret provides a framework to create tools to inspect information
looking for sensitive information like passwords, tokens, private keys,
certificates, etc.



Basics

The current trend of automation of all things and de DevOps culture are very 
beneficial for efficiency but also come with several  problems, being one of
them the secret provisioning. Bootstrapping secrets into systems and
applications may be complicated and sometimes the straightforward way is to 
store them into a insecure storage, like github repository, embedded into an
artifact or system image, etc. That means that an AWS secret_key end up into a
Github repository.

Seekret is an extensible framework that gelps in creating tools for detecting
secrets on different sources. The secrets to detect are defined by a set of
rules that can help detect passwords, tokens, private keys, certificates, etc.



Tools Using Seekret

Seekret is extensible and can cover various use cases. Below there are some
tools that uses seekret:

  git-seekret: https://github.com/apuigsech/git-seekret

    Git module that uses local hooks to help develepers to prevent leaking 
    sensitive information in a commit.



Using It

Seekret API is very simple and easy to use. This section shows some snippets of
code that shows the basic operations you can do with it.

The first thing to be done is to create a new Seekret context:

	s := seekret.NewSeekret()

Then the rules must to be loaded. They can be loaded from a path definition, a
directory or a single file:

	s.LoadRulesFromPath("/path/to/main/rues:/path/to/other/rules:/path/to/more/rules")

	s.LoadRulesFromDir("/path/to/rules")

	s.LoadRulesFromFile("/path/to/file.rule")

Optionally, exceptions (or false positives) can also be loaded from a file:

	s.LoadExceptionsFromFile("/path/to/exceptions/file")

After that, must be loaded the objects to be inspected searching for secrets.

	opts := map[string]interface{} {
  		// Loading options.
	}
	s.LoadObjects(sourceType, source, opts)

sourceType is an interface that implements the interface shown below. We offer
sourceType's for Directories and Git Repositories, but you are able to extend
it by creating your own.

	type SourceType interface {
		LoadObjects(source string, opt LoadOptions) ([]models.Object, error)
	}

Currently, there are the following different sources supported:

  Directories (and files): https://github.com/apuigsech/seekret-source-dir

    Load all files contained in a directory (and its sub-directories).

  Git Repositories: https://github.com/apuigsech/seekret-source-git

    Load git objects from commits or staging area.

Having all the rules, exceptions and objects loaded into the contects, it's
possible to start the inspection with the following code:

	s.Inspect(Nworkers)

Nworkers is an integuer that specify the number of goroutines used on the
inspection. The recommended value is runtime.NumCPU().

Finally, it is possible to obtain the list of secrets located and do
something with them:

		secretsList := s.ListSecrets()
		for secret := range secretsList {
			// Do something
		}

*/
package seekret
